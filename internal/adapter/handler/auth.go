package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type signupRequest struct {
	Name     string `json:"name" validate:"required,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12,max=128"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type authResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type AuthHandler struct {
	service    ports.AuthServicer
	cookieName string
}

func NewAuthHandler(as ports.AuthServicer) *AuthHandler {
	return &AuthHandler{
		service:    as,
		cookieName: "refresh_token",
	}
}

func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	signupCmd := &ports.SignupCmd{
		Name: req.Name,
		LoginCmd: ports.LoginCmd{
			Email:     req.Email,
			Password:  req.Password,
			UserAgent: r.UserAgent(),
		},
	}

	res, err := ah.service.Signup(r.Context(), signupCmd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	ah.setRefreshToken(w, res.RefreshToken, time.Now().Add(ah.service.RefreshTokenTTL()))
	writeJson(w, http.StatusOK, authResponse{
		AccessToken: res.AccessToken,
		ExpiresAt:   res.ExpiresAt,
	})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid json"))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}

	loginCmd := &ports.LoginCmd{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: r.UserAgent(),
	}
	res, err := ah.service.Login(r.Context(), loginCmd)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}

	ah.setRefreshToken(w, res.RefreshToken, time.Now().Add(ah.service.RefreshTokenTTL()))
	writeJson(w, http.StatusOK, authResponse{
		AccessToken: res.AccessToken,
		ExpiresAt:   res.ExpiresAt,
	})
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	rt, ok := ah.getRefreshToken(r)
	if ok {
		_ = ah.service.Logout(r.Context(), &ports.LogoutCmd{
			Token: rt,
		})
	}

	ah.clearRefreshToken(w)
	w.WriteHeader(http.StatusNoContent)
}

func (ah *AuthHandler) setRefreshToken(w http.ResponseWriter, token string, expires time.Time) {
	c := &http.Cookie{
		Name:     ah.cookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
}

func (ah *AuthHandler) clearRefreshToken(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     ah.cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
}

func (ah *AuthHandler) getRefreshToken(r *http.Request) (string, bool) {
	c, err := r.Cookie(ah.cookieName)
	if err != nil || c.Value == "" {
		return "", false
	}

	return c.Value, true
}
