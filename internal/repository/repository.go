package repository

import "github.com/zapkub/cftl/internal/database"

type DB struct {
	db *database.DB
}

func New(db *database.DB) *DB {
	return &DB{db: db}
}
