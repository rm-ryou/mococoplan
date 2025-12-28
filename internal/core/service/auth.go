package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"time"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
	"github.com/rm-ryou/mococoplan/pkg/password"
)

type AuthService struct {
	issuer      ports.TokenIssuer
	userRepo    ports.UserRepository
	sessionRepo ports.SessionRepository

	params          *password.Params
	refreshTokenTTL time.Duration
}

func NewAuthService(
	issuer ports.TokenIssuer,
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	params *password.Params,
	refreshTokenTTL time.Duration,
) ports.AuthServicer {
	return &AuthService{
		issuer:          issuer,
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		params:          params,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (as *AuthService) Signup(ctx context.Context, cmd *ports.SignupCmd) (*ports.AuthResult, error) {
	hash, err := password.Hash(cmd.Password, as.params)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Name:         cmd.Name,
		Email:        cmd.Email,
		PasswordHash: hash,
	}

	if err := as.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	newUser, err := as.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	return as.issue(ctx, newUser, cmd.IP, cmd.UserAgent)
}

func (as *AuthService) Login(ctx context.Context, cmd *ports.LoginCmd) (*ports.AuthResult, error) {
	u, err := as.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}

	ok, err := password.Verify(cmd.Password, u.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrInvalidCredentials
	}

	return as.issue(ctx, u, cmd.IP, cmd.UserAgent)
}

func (as *AuthService) Logout(ctx context.Context, cmd *ports.LogoutCmd) error {
	hash := sha256.Sum256([]byte(cmd.Token))
	err := as.sessionRepo.Delete(ctx, hash)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) RefreshTokenTTL() time.Duration {
	return as.refreshTokenTTL
}

func (as *AuthService) issue(ctx context.Context, u *domain.User, ip net.IP, ua string) (*ports.AuthResult, error) {
	identity := &ports.UserIdentity{
		UserID: u.ID,
		Email:  u.Email,
	}

	at, err := as.issuer.Issue(identity)
	if err != nil {
		return nil, err
	}

	// FIXME: move generate refresh token logic
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	rt := base64.RawURLEncoding.EncodeToString(b)
	rth := sha256.Sum256([]byte(rt))

	s := &domain.Session{
		UserID:    u.ID,
		Token:     rth,
		IP:        toIP16(ip),
		UserAgent: ua,
		ExpiresAt: time.Now().Add(as.RefreshTokenTTL()),
	}
	if err := as.sessionRepo.Create(ctx, s); err != nil {
		return nil, err
	}

	return &ports.AuthResult{
		AccessToken:  at.Token,
		ExpiresAt:    at.ExpiresAt,
		RefreshToken: rt,
	}, nil
}

func toIP16(ip net.IP) [16]byte {
	var res [16]byte
	if ip == nil {
		return res
	}

	if v4 := ip.To4(); v4 != nil {
		res[10], res[11] = 0xff, 0xff
		copy(res[12:], v4)
		return res
	}

	if v6 := ip.To16(); v6 != nil {
		copy(res[:], v6)
	}
	return res
}
