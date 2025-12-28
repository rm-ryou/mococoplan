package jwt

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
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

func (s *Service) Issue(identity *ports.UserIdentity) (*ports.AccessToken, error) {
	now := time.Now()
	exp := now.Add(s.ttl)

	jc := &jwtClaims{
		UserId: strconv.Itoa(identity.UserID),
		Email:  identity.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   strconv.Itoa(identity.UserID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jc)
	ss, err := t.SignedString(s.secret)
	if err != nil {
		return nil, err
	}

	return &ports.AccessToken{
		Token:     ss,
		ExpiresAt: exp,
	}, nil
}

func (s *Service) Verify(t string) (*ports.UserIdentity, error) {
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

	userID, err := strconv.Atoi(jc.UserId)
	if err != nil {
		return nil, err
	}

	return &ports.UserIdentity{
		UserID: userID,
		Email:  jc.Email,
	}, nil
}
