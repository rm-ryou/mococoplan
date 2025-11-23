package router

import (
	"database/sql"
	"net/http"
)

func Setup(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	NewUserRouter(mux, db)

	return mux
}
