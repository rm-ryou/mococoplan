package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rm-ryou/mococoplan/internal/core/token"
)

var ErrInvalidToken = errors.New("invalid token")

type Service struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type jwtClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func New(secret, issuer string, ttl time.Duration) *Service {
	return &Service{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (s *Service) Issue(c *token.Claims) (*token.AccessToken, error) {
	now := time.Now()
	exp := now.Add(s.ttl)

	jc := &jwtClaims{
		UserId: strconv.Itoa(c.UserId),
		Email:  c.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   strconv.Itoa(c.UserId),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jc)
	ss, err := t.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &token.AccessToken{
		Token:     ss,
		ExpiresAt: exp,
	}, nil
}

func (s *Service) Verify(t string) (*token.Claims, error) {
	jc := &jwtClaims{}
	parsed, err := jwt.ParseWithClaims(t, jc, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})

	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	if jc.UserId == "" || jc.Email == "" {
		return nil, ErrInvalidToken
	}

	userId, err := strconv.Atoi(jc.UserId)
	if err != nil {
		return nil, err
	}

	return &token.Claims{
		UserId: userId,
		Email:  jc.Email,
	}, nil
}
