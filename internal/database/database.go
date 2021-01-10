package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
	mu sync.Mutex

	tx *sql.Tx
}

func Open(driverName, dbinfo string) (_ *DB, err error) {

	db, err := sql.Open(driverName, dbinfo)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return New(db), nil
}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (_ *sql.Rows, err error) {
	return db.db.QueryContext(ctx, query, args...)
}
func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

// Exec executes a SQL statement and returns the number of rows it affected.
func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (_ int64, err error) {
	res, err := db.execResult(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("RowsAffected: %v", err)
	}
	return n, nil
}

// execResult executes a SQL statement and returns a sql.Result.
func (db *DB) execResult(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	return db.db.ExecContext(ctx, query, args...)
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Transact(ctx context.Context, iso sql.IsolationLevel, txfn func(*DB) error) (err error) {

	if db.InTransaction() {
		return errors.New("a DB transaction has been called in transactioned db")
	}

	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("db.BeginTx(): %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			if txErr := tx.Commit(); txErr != nil {
				err = fmt.Errorf("txCommit(): %w", txErr)
			}
		}
	}()

	dbtx := New(db.db)
	dbtx.tx = tx
	if err := txfn(dbtx); err != nil {
		return fmt.Errorf("txfn(tx) return error: %w", err)
	}
	return nil
}

func (db *DB) InTransaction() bool {
	return db.tx != nil
}
