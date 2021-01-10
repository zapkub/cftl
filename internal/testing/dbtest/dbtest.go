package dbtest

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/zapkub/cftl/internal/apperror"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// DBConnURI generates a postgres connection string in URI format.  This is
// necessary as migrate expects a URI.
func DBConnURI(dbName string) string {
	var (
		user     = getEnv("CFTL_DATABASE_TEST_USER", "postgres")
		password = getEnv("CFTL_DATABASE_TEST_PASSWORD", "root")
		host     = getEnv("CFTL_DATABASE_TEST_HOST", "localhost")
		port     = getEnv("CFTL_DATABASE_TEST_PORT", "5432")
	)
	cs := fmt.Sprintf("postgres://%s/%s?sslmode=disable&user=%s&password=%s&port=%s",
		host, dbName, url.QueryEscape(user), url.QueryEscape(password), url.QueryEscape(port))
	return cs
}

// CreateDB creates a new database dbName.
func CreateDB(dbName string) error {
	return ConnectAndExecute(DBConnURI(""), func(pg *sql.DB) error {
		if _, err := pg.Exec(fmt.Sprintf(`
			CREATE DATABASE %q
				TEMPLATE=template0
				LC_COLLATE='C'
				LC_CTYPE='C';`, dbName)); err != nil {
			return fmt.Errorf("error creating %q: %v", dbName, err)
		}

		return nil
	})
}

// ConnectAndExecute connects to the postgres database specified by uri and
// executes dbFunc, then cleans up the database connection.
// It returns an error that Is derrors.NotFound if no connection could be made.
func ConnectAndExecute(uri string, dbFunc func(*sql.DB) error) (outerErr error) {
	pg, err := sql.Open("postgres", uri)
	if err == nil {
		err = pg.Ping()
	}
	if err != nil {
		return fmt.Errorf("%w: %v", apperror.NotFound, err)
	}
	defer func() {
		if err := pg.Close(); err != nil {
			outerErr = MultiErr{outerErr, err}
		}
	}()
	return dbFunc(pg)
}

// MultiErr can be used to combine one or more errors into a single error.
type MultiErr []error

func (m MultiErr) Error() string {
	var sb strings.Builder
	for _, err := range m {
		sep := ""
		if sb.Len() > 0 {
			sep = "|"
		}
		if err != nil {
			sb.WriteString(sep + err.Error())
		}
	}
	return sb.String()
}

// checkIfDBExists check if dbName exists.
func checkIfDBExists(dbName string) (bool, error) {
	var exists bool

	err := ConnectAndExecute(DBConnURI(""), func(pg *sql.DB) error {
		rows, err := pg.Query("SELECT 1 from pg_database WHERE datname = $1 LIMIT 1", dbName)
		if err != nil {
			return err
		}
		defer rows.Close()

		if rows.Next() {
			exists = true
			return nil
		}

		return rows.Err()
	})

	return exists, err
}

// CreateDBIfNotExists checks whether the given dbName is an existing database,
// and creates one if not.
func CreateDBIfNotExists(dbName string) error {
	exists, err := checkIfDBExists(dbName)
	if err != nil || exists {
		return err
	}

	log.Printf("Test database %q does not exist, creating.", dbName)
	return CreateDB(dbName)
}

// DropDB drops the database named dbName.
func DropDB(dbName string) error {
	return ConnectAndExecute(DBConnURI(""), func(pg *sql.DB) error {
		if _, err := pg.Exec(fmt.Sprintf("DROP DATABASE %q;", dbName)); err != nil {
			return fmt.Errorf("error dropping %q: %v", dbName, err)
		}
		return nil
	})
}
