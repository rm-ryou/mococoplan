package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/rm-ryou/mococoplan/internal/core/ports"
)

type ctxKey string

const userIdentityKey ctxKey = "user-identity"

type Auth struct {
	Verifier ports.TokenVerifier
}

func NewAuth(verifier ports.TokenVerifier) *Auth {
	return &Auth{
		Verifier: verifier,
	}
}

func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		ui, err := a.Verifier.Verify(token)
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIdentityKey, ui)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdentityFromCtx(ctx context.Context) (ports.UserIdentity, bool) {
	val := ctx.Value(userIdentityKey)
	if val == nil {
		return ports.UserIdentity{}, false
	}

	p, ok := val.(ports.UserIdentity)
	return p, ok
}

func bearerToken(val string) (string, error) {
	if val == "" {
		return "", errors.New("authorization header not setted")
	}
	parts := strings.SplitN(val, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid token format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
