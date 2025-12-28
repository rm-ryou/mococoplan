package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
)

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

func TestUserRepository_SuccessCreate(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	repo := NewUserRepository(testDB)

	user := &domain.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	records := countUsersByEmail(t, testDB, user.Email)
	if records != 1 {
		t.Errorf("want user count: %d, act: %d", 1, records)
	}
}

func TestUserRepository_FailedCreate_DuplicateEmail(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	repo := NewUserRepository(testDB)

	ctx := context.Background()
	user := &domain.User{
		Name:         "test name",
		Email:        "dup@example.com",
		PasswordHash: "testHashedPassword",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Create(ctx, user)
	if err != domain.ErrEmailAlreadyExists {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUserRepository_SuccessFindByEmail(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	repo := NewUserRepository(testDB)

	user := &domain.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	act, err := repo.FindByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if act.Email != user.Email {
		t.Errorf("want: %v, act: %v", user.Email, act.Email)
	}
}

func TestUserRepository_FailedFindByEmail_RecordNotExists(t *testing.T) {
	defer deleteAllRecords(t, testDB, "users")
	repo := NewUserRepository(testDB)

	email := "not-exists@example.com"

	act, err := repo.FindByEmail(context.Background(), email)
	if err != domain.ErrUserNotFound {
		t.Fatalf("unexpected error: %v", err)
	}

	if act != nil {
		t.Fatalf("want: %v, act: %v", nil, act)
	}
}
