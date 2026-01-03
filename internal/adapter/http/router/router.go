package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/config"
)

func Setup(db *sql.DB, tokenCfg config.Token) http.Handler {
	mux := http.NewServeMux()

	NewAuthRouter(mux, db, tokenCfg)
	NewWorkspaceRouter(mux, db, tokenCfg)

	return mux
}
