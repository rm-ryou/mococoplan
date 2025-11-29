package handler

import (
	"io"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/core/port"
)

type UserHandler struct {
	us port.UserServicer
}

func NewUserHandler(us port.UserServicer) *UserHandler {
	return &UserHandler{us: us}
}

func (c *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is SignUp page!!")
}

func (c *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is SignIn page!!")
}

func (c *UserHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is SignOut page!!")
}

func (c *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is VerifyEmail page!!")
}
