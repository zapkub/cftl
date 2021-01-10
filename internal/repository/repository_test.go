package repository

import "testing"

var testDB *DB

func TestMain(m *testing.M) {
	RunDBTests("cftl_repository_test", m, &testDB)
}
