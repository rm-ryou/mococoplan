package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/adapter/handler"
	"github.com/rm-ryou/mococoplan/internal/adapter/repository"
	"github.com/rm-ryou/mococoplan/internal/adapter/token/jwt"
	"github.com/rm-ryou/mococoplan/internal/config"
	"github.com/rm-ryou/mococoplan/internal/core/service"
	"github.com/rm-ryou/mococoplan/pkg/password"
)

func NewAuthRouter(mux *http.ServeMux, db *sql.DB, tokenCfg config.Token) {
	issuer := jwt.New(tokenCfg.AccessTokenSecret, tokenCfg.AccessTokenSecret, tokenCfg.AccessTokenTTL)
	ur := repository.NewUserRepository(db)
	sr := repository.NewSessionRepository(db)
	pp := password.DefaultParams()

	as := service.NewAuthService(issuer, ur, sr, pp, tokenCfg.RefreshTokenTTL)
	h := handler.NewAuthHandler(as)

	authMux := http.NewServeMux()
	authMux.HandleFunc("POST /api/v1/auth/signup", h.Signup)
	authMux.HandleFunc("POST /api/v1/auth/login", h.Login)
	authMux.HandleFunc("POST /api/v1/auth/logout", h.Logout)

	mux.Handle("/api/v1/auth/", authMux)
}
