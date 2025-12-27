package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"testing"
	"time"

	"github.com/rm-ryou/mococoplan/internal/core/domain/session"
	"github.com/rm-ryou/mococoplan/internal/core/domain/user"
)

func countSessionsRecords(t *testing.T, db *sql.DB, token [32]byte) int {
	t.Helper()

	query := "SELECT COUNT(*) FROM sessions WHERE token = ?"
	row := db.QueryRow(query, token[:])

	var cnt int
	if err := row.Scan(&cnt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return cnt
}

func TestSessionRepository_SuccessCreate(t *testing.T) {
	defer deleteAllRecords(t, testDB, "sessions")
	defer deleteAllRecords(t, testDB, "users")
	userRepo := NewUserRepository(testDB)
	repo := NewSessionRepository(testDB)

	token := sha256.Sum256([]byte("test-token"))

	if err := userRepo.Create(context.Background(), &user.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	testSession := &session.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now(),
	}

	err = repo.Create(context.Background(), testSession)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	records := countSessionsRecords(t, testDB, token)
	if records != 1 {
		t.Errorf("want session count: %d, act: %d", 1, records)
	}
}

func TestSessionRepository_FailedCreate_DuplicateSession(t *testing.T) {
	defer deleteAllRecords(t, testDB, "sessions")
	defer deleteAllRecords(t, testDB, "users")
	userRepo := NewUserRepository(testDB)
	repo := NewSessionRepository(testDB)

	token := sha256.Sum256([]byte("test-token"))

	if err := userRepo.Create(context.Background(), &user.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	testSession := &session.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now(),
	}

	err = repo.Create(context.Background(), testSession)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = repo.Create(context.Background(), testSession)
	if err != session.ErrInvalid {
		t.Errorf("want: %v, act: %v", session.ErrInvalid, err)
	}
}

func TestSessionRepository_FindByToken(t *testing.T) {
	defer deleteAllRecords(t, testDB, "sessions")
	defer deleteAllRecords(t, testDB, "users")
	userRepo := NewUserRepository(testDB)
	repo := NewSessionRepository(testDB)

	token := sha256.Sum256([]byte("test-token"))

	if err := userRepo.Create(context.Background(), &user.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	testSession := &session.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now(),
	}

	err = repo.Create(context.Background(), testSession)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("Success to Get session", func(t *testing.T) {
		s, err := repo.FindByToken(context.Background(), testSession.Token)

		if s == nil {
			t.Errorf("want: %v, act: %v", nil, s)
		}
		if err != nil {
			t.Errorf("want: %v, act: %v", nil, err)
		}
	})

	t.Run("When token is not saved, return session.ErrNotFound error", func(t *testing.T) {
		s, err := repo.FindByToken(context.Background(), sha256.Sum256([]byte("no-exists-token")))

		if s != nil {
			t.Errorf("want: %v, act: %v", nil, s)
		}
		if err != session.ErrNotFound {
			t.Errorf("want: %v, act: %v", session.ErrNotFound, err)
		}
	})
}

func TestSessionRepository_SuccessDelete(t *testing.T) {
	defer deleteAllRecords(t, testDB, "sessions")
	defer deleteAllRecords(t, testDB, "users")
	userRepo := NewUserRepository(testDB)
	repo := NewSessionRepository(testDB)

	token := sha256.Sum256([]byte("test-token"))

	if err := userRepo.Create(context.Background(), &user.User{
		Name:         "test name",
		Email:        "test@example.com",
		PasswordHash: "testHashedPassword",
	}); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := userRepo.FindByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	testSession := &session.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now(),
	}

	err = repo.Create(context.Background(), testSession)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Delete(context.Background(), testSession.Token)

	if err != nil {
		t.Errorf("want: %v, act: %v", nil, err)
	}

	sessionRecords := countSessionsRecords(t, testDB, testSession.Token)
	if sessionRecords != 0 {
		t.Errorf("records should be 0, want: %v, act: %v", 0, sessionRecords)
	}

}

func TestSessionRepository_SuccessDelete_WhenSessionNotSaved(t *testing.T) {
	defer deleteAllRecords(t, testDB, "sessions")
	repo := NewSessionRepository(testDB)

	token := sha256.Sum256([]byte("no-exists-token"))

	err := repo.Delete(context.Background(), token)

	if err != nil {
		t.Errorf("want: %v, act: %v", nil, err)
	}
}
