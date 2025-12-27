package repository

import (
	"database/sql"
	"testing"
)

func deleteAllRecords(t *testing.T, db *sql.DB, tableName string) {
	t.Helper()

	if _, err := db.Exec("DELETE FROM " + tableName); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
