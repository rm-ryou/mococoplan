package router

import (
	"database/sql"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/adapter/http/handler"
	"github.com/rm-ryou/mococoplan/internal/adapter/http/middleware"
	"github.com/rm-ryou/mococoplan/internal/adapter/repository"
	"github.com/rm-ryou/mococoplan/internal/adapter/token/jwt"
	"github.com/rm-ryou/mococoplan/internal/config"
	"github.com/rm-ryou/mococoplan/internal/core/service"
)

func NewWorkspaceRouter(mux *http.ServeMux, db *sql.DB, tokenCfg config.Token) {
	tx := repository.NewTxManager(db)
	wsr := repository.NewWorkspaceRepository(db)
	wsmr := repository.NewWorkspaceMemberRepository(db)

	ws := service.NewWorkspaceService(tx, wsr, wsmr)
	h := handler.NewWorkspaceHandler(ws)

	verifier := jwt.New(tokenCfg.AccessTokenSecret, tokenCfg.AccessTokenSecret, tokenCfg.AccessTokenTTL)
	auth := middleware.NewAuth(verifier)

	workspaceMux := http.NewServeMux()
	workspaceMux.Handle("POST /api/v1/workspaces/", auth.RequireAuth(http.HandlerFunc(h.Create)))
	workspaceMux.Handle("GET /api/v1/workspaces/", auth.RequireAuth(http.HandlerFunc(h.ListWorkspaces)))
	workspaceMux.Handle("POST /api/v1/workspaces/{workspacesId}/members", auth.RequireAuth(http.HandlerFunc(h.AddMember)))
	workspaceMux.Handle("GET /api/v1/workspaces/{workspacesId}/members", auth.RequireAuth(http.HandlerFunc(h.ListMembers)))

	mux.Handle("/api/v1/workspaces/", workspaceMux)
}
