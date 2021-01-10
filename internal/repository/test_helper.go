package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/zapkub/cftl/internal/apperror"
	"github.com/zapkub/cftl/internal/database"
	"github.com/zapkub/cftl/internal/testing/dbtest"

	// imported to register the postgres migration driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// imported to register the file source migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// imported to register the postgres database driver
	_ "github.com/lib/pq"
)

// RunDBTests is a wrapper that runs the given testing suite in a test database
// named dbName.  The given *DB reference will be set to the instantiated test
// database.
func RunDBTests(dbName string, m *testing.M, testDB **DB) {
	if err := dbtest.CreateDBIfNotExists(dbName); err != nil {
		log.Fatal(err)
	}
	db, err := SetupTestDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		if errors.Is(err, apperror.NotFound) && os.Getenv("GO_DISCOVERY_TESTDB") != "true" {
			log.Printf("SKIPPING: could not connect to DB (see doc/postgres.md to set up): %v", err)
			return
		}
		log.Fatal(err)
	}
	*testDB = db
	code := m.Run()
	if err := db.db.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

// SetupTestDB creates a test database named dbName if it does not already
// exist, and migrates it to the latest schema from the migrations directory.
func SetupTestDB(dbName string) (_ *DB, err error) {

	if err := dbtest.CreateDBIfNotExists(dbName); err != nil {
		return nil, fmt.Errorf("CreateDBIfNotExists(%q): %w", dbName, err)
	}
	if isMigrationError, err := tryToMigrate(dbName); err != nil {
		if isMigrationError {
			// failed during migration stage, recreate and try again
			log.Printf("Migration failed for %s: %v, recreating database.", dbName, err)
			if err := recreateDB(dbName); err != nil {
				return nil, fmt.Errorf("recreateDB(%q): %v", dbName, err)
			}
			_, err = tryToMigrate(dbName)
		}
		if err != nil {
			return nil, fmt.Errorf("unfixable error migrating database: %v", err)
		}
	}
	driver := os.Getenv("GO_DISCOVERY_DATABASE_DRIVER")
	if driver == "" {
		driver = "postgres"
	}
	db, err := database.Open(driver, dbtest.DBConnURI(dbName))
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

// recreateDB drops and recreates the database named dbName.
func recreateDB(dbName string) error {
	err := dbtest.DropDB(dbName)
	if err != nil {
		return err
	}

	return dbtest.CreateDB(dbName)
}

// tryToMigrate attempts to migrate the database named dbName to the latest
// migration. If this operation fails in the migration step, it returns
// isMigrationError=true to signal that the database should be recreated.
func tryToMigrate(dbName string) (isMigrationError bool, outerErr error) {
	dbURI := dbtest.DBConnURI(dbName)
	wd, _ := os.Getwd()
	source := path.Join(wd, "../../migrations")
	m, err := migrate.New("file://"+source, dbURI)
	if err != nil {
		return false, fmt.Errorf("migrate.New(): %v", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			outerErr = dbtest.MultiErr{outerErr, srcErr, dbErr}
		}
	}()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return true, fmt.Errorf("m.Up(): %v", err)
	}
	return false, nil
}

func ResetTestDB(db *DB, t *testing.T) {

	ctx := context.Background()
	t.Helper()

	if err := db.db.Transact(ctx, sql.LevelDefault, func(d *database.DB) error {
		if _, err := d.Exec(ctx, `
			TRUNCATE users CASCADE;
			TRUNCATE sessions CASCADE;
			TRUNCATE github_oauths CASCADE;
		`); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatalf("resetTestDB(): %+v", err)
	}

}
