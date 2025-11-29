package user

import "context"

type Creater interface {
	Create(ctx context.Context, name, email, password string) error
}
