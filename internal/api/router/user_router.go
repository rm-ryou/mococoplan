package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/api/controller"
	"github.com/rm-ryou/mococoplan/internal/repository"
	"github.com/rm-ryou/mococoplan/internal/service"
)

func NewUserRouter(mux *http.ServeMux, db *sql.DB) {
	r := repository.NewUserRepository(db)
	s := service.NewUserService(r)
	c := controller.NewUserController(s)

	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /api/v1/users/sign_up", c.SignUp)

	mux.Handle("/api/v1/users/", userMux)
}
