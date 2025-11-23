package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/rm-ryou/mococoplan/internal/domain"
)

type UserServicer interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
}

type UserController struct {
	us UserServicer
}

func NewUserController(us UserServicer) *UserController {
	return &UserController{us: us}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is SignUp page!!")
}
