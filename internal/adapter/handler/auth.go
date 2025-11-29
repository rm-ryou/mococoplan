package handler

import (
	"io"
	"net/http"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is signup page")
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is login page")
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is logout page")
}
