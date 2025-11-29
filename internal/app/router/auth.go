package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/adapter/handler"
)

func NewAuthRouter(mux *http.ServeMux, db *sql.DB) {
	h := handler.NewAuthHandler()

	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /api/v1/auth/signup", h.Signup)
	authMux.HandleFunc("POST /api/v1/auth/login", h.Login)
	authMux.HandleFunc("POST /api/v1/auth/logout", h.Logout)

	mux.Handle("/api/v1/auth/", authMux)
}
