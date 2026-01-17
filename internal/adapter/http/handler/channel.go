package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/rm-ryou/mococoplan/internal/adapter/http/middleware"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type createChannelRequest struct {
	WorkspaceID int    `json:"workspace_id" validate:"required"`
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description" validate:"required,max=255"`
	IsPrivate   bool   `json:"is_private"`
}

type ChannelHandler struct {
	service ports.ChannelServicer
}

func NewChannelHandler(cs ports.ChannelServicer) *ChannelHandler {
	return &ChannelHandler{
		service: cs,
	}
}

func (ch *ChannelHandler) Create(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	var req createChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	cmd := &ports.CreateChannelCmd{
		UserID:      identity.UserID,
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	channel, err := ch.service.Create(r.Context(), cmd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusOK, channel)
}

func (ch *ChannelHandler) ListJoined(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	wsIDStr := r.PathValue("workspaceID")
	wsID, err := strconv.Atoi(wsIDStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errors.New("invalid workspace"))
	}

	items, err := ch.service.ListJoined(r.Context(), wsID, identity.UserID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
	}

	writeJson(w, http.StatusOK, items)
}

func (ch *ChannelHandler) Join(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	chanIDStr := r.PathValue("channelID")
	chanID, err := strconv.Atoi(chanIDStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errors.New("invalid workspace"))
	}

	if err := ch.service.Join(r.Context(), chanID, identity.UserID); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChannelHandler) Leave(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	chanIDStr := r.PathValue("channelID")
	chanID, err := strconv.Atoi(chanIDStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errors.New("invalid workspace"))
	}

	if err := ch.service.Leave(r.Context(), chanID, identity.UserID); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ch *ChannelHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	identity, ok := middleware.UserIdentityFromCtx(r.Context())
	if !ok {
		writeErr(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	chanIDStr := r.PathValue("channelID")
	chanID, err := strconv.Atoi(chanIDStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, errors.New("invalid workspace"))
	}

	if err := ch.service.AddMember(r.Context(), chanID, identity.UserID); err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
