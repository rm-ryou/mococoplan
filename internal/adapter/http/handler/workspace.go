package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rm-ryou/mococoplan/internal/adapter/http/middleware"
	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type createWorkspaceRequest struct {
	Name string `json:"name" validate:"required,max=50"`
	Slug string `json:"slug" validate:"required,max=50"`
}

type addWorkspaceMemberRequest struct {
	UserID int    `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required"`
}

type WorkspaceHandler struct {
	service ports.WorkspaceServicer
}

func NewWorkspaceHandler(ws ports.WorkspaceServicer) *WorkspaceHandler {
	return &WorkspaceHandler{
		service: ws,
	}
}

func (wh *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	cmd := &ports.CreateWorkspaceCmd{
		UserID: identity.UserID,
		Name:   req.Name,
		Slug:   req.Slug,
	}

	ws, err := wh.service.Create(r.Context(), cmd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]any{
		"id":   ws.ID,
		"name": ws.Name,
		"slug": ws.Slug,
	})
}

func (wh *WorkspaceHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var req addWorkspaceMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	workspaceIDStr := r.PathValue("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	cmd := &ports.AddWorkspaceMemberCmd{
		UserID:       identity.UserID,
		WorkspaceID:  workspaceID,
		TargetUserID: req.UserID,
		Role:         domain.WorkspaceRole(req.Role),
	}

	err = wh.service.AddMember(r.Context(), cmd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (wh *WorkspaceHandler) ListWorkspaces(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	items, err := wh.service.ListWorkspaces(r.Context(), identity.UserID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	workspaces := make([]map[string]any, 0, len(items))
	for _, ws := range items {
		workspaces = append(workspaces, map[string]any{
			"id":   ws.ID,
			"name": ws.Name,
			"slug": ws.Slug,
		})
	}
	writeJson(w, http.StatusOK, map[string]any{
		"workspaces": workspaces,
	})
}

func (wh *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	workspaceIDStr := r.PathValue("workspaceID")
	workspaceID, err := strconv.Atoi(workspaceIDStr)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	items, err := wh.service.ListMembers(r.Context(), workspaceID, identity.UserID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	members := make([]map[string]any, 0, len(items))
	for _, m := range items {
		members = append(members, map[string]any{
			"user_id":   m.UserID,
			"role":      m.Role,
			"joined_at": m.JoinedAt,
		})
	}
}
