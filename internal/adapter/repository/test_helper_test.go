package repository

import (
	"database/sql"
	"testing"
)

func deleteAllUsers(t *testing.T, db *sql.DB) {
	t.Helper()

	if _, err := db.Exec("DELETE FROM users"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func countUsersByEmail(t *testing.T, db *sql.DB, email string) int {
	t.Helper()

	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	row := db.QueryRow(query, email)

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return cnt
}
