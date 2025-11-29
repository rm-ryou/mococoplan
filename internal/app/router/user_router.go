package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/adapter/handler"
	"github.com/rm-ryou/mococoplan/internal/adapter/repository"
	"github.com/rm-ryou/mococoplan/internal/core/service"
)

func NewUserRouter(mux *http.ServeMux, db *sql.DB) {
	r := repository.NewUserRepository(db)
	s := service.NewUserService(r)
	c := handler.NewUserController(s)

	userMux := http.NewServeMux()
	userMux.HandleFunc("POST /api/v1/users/sign_up", c.SignUp)

	mux.Handle("/api/v1/users/", userMux)
}
