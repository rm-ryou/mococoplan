package repository

import (
	"database/sql"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	host := getEnv("TEST_DB_HOST", "127.0.0.1")
	port := getEnv("TEST_DB_PORT", "3306")
	name := getEnv("TEST_DB_NAME", "mococoplan")
	user := getEnv("TEST_DB_USER", "user")
	password := getEnv("TEST_DB_PASSWORD", "password")

	dsn := CreateDSN(name, user, password, host, port)
	var err error
	testDB, err = NewDB(dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	_ = testDB.Close()
	os.Exit(code)
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
