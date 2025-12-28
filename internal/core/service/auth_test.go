package service

import (
	"context"
	"crypto/sha256"
	"net"
	"testing"
	"time"

	"github.com/rm-ryou/mococoplan/internal/core/domain"
	"github.com/rm-ryou/mococoplan/internal/core/ports"
	"github.com/rm-ryou/mococoplan/internal/core/ports/mocks"
	"github.com/rm-ryou/mococoplan/pkg/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_SuccessSignup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	issuer := new(mocks.TokenIssuer)
	userRepo := new(mocks.UserRepository)
	sessionRepo := new(mocks.SessionRepository)
	pp := password.DefaultParams()
	refreshTokenTTL := 30 * 24 * time.Hour

	service := NewAuthService(issuer, userRepo, sessionRepo, pp, refreshTokenTTL)

	cmd := &ports.SignupCmd{
		Name: "Test",
		LoginCmd: ports.LoginCmd{
			Email:     "test@example.com",
			Password:  "test-password",
			IP:        net.ParseIP("127.0.0.1"),
			UserAgent: "TestUserAgent",
		},
	}

	userRepo.On("Create", ctx, mock.MatchedBy(func(u *domain.User) bool {
		return u.Name == cmd.Name &&
			u.Email == cmd.Email &&
			len(u.PasswordHash) > 0
	})).Return(nil).Once()

	createdUser := &domain.User{
		ID:           1,
		Name:         cmd.Name,
		Email:        cmd.Email,
		PasswordHash: "hashedPassword",
	}
	userRepo.On("FindByEmail", ctx, cmd.Email).Return(createdUser, nil).Once()

	issued := &ports.AccessToken{
		Token:     "access.jwt",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	issuer.On("Issue", mock.Anything).Return(issued, nil).Once()

	sessionRepo.On("Create", ctx, mock.MatchedBy(func(s *domain.Session) bool {
		if s.UserID != createdUser.ID {
			return false
		}
		if s.UserAgent != cmd.UserAgent {
			return false
		}
		return len(s.Token) == 32
	})).Return(nil).Once()

	res, err := service.Signup(ctx, cmd)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, issued.Token, res.AccessToken)
	assert.Equal(t, issued.ExpiresAt, res.ExpiresAt)
	assert.NotEmpty(t, res.RefreshToken)

	issuer.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestAuthService_SuccessLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	issuer := new(mocks.TokenIssuer)
	userRepo := new(mocks.UserRepository)
	sessionRepo := new(mocks.SessionRepository)
	pp := password.DefaultParams()
	refreshTokenTTL := 30 * 24 * time.Hour

	service := NewAuthService(issuer, userRepo, sessionRepo, pp, refreshTokenTTL)

	testPassword := "test-password"
	hash, err := password.Hash(testPassword, pp)
	require.NoError(t, err)

	u := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hash,
	}

	cmd := &ports.LoginCmd{
		Email:     "test@example.com",
		Password:  testPassword,
		IP:        net.ParseIP("127.0.0.1"),
		UserAgent: "TestUserAgent",
	}

	userRepo.On("FindByEmail", ctx, cmd.Email).Return(u, nil).Once()
	issued := &ports.AccessToken{
		Token:     "access.jwt",
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	issuer.On("Issue", mock.Anything).Return(issued, nil).Once()

	sessionRepo.On("Create", ctx, mock.MatchedBy(func(s *domain.Session) bool {
		return s.UserID == u.ID && s.UserAgent == cmd.UserAgent
	})).Return(nil).Once()

	res, err := service.Login(ctx, cmd)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, issued.Token, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)

	issuer.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestAuthService_FailedLogin_InvalidCredentials(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	issuer := new(mocks.TokenIssuer)
	userRepo := new(mocks.UserRepository)
	sessionRepo := new(mocks.SessionRepository)
	pp := password.DefaultParams()
	refreshTokenTTL := 30 * 24 * time.Hour

	service := NewAuthService(issuer, userRepo, sessionRepo, pp, refreshTokenTTL)

	testPassword := "test-password"
	hash, err := password.Hash(testPassword, pp)
	require.NoError(t, err)

	u := &domain.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: hash,
	}

	cmd := &ports.LoginCmd{
		Email:     "test@example.com",
		Password:  "test-wrong-password",
		IP:        net.ParseIP("127.0.0.1"),
		UserAgent: "TestUserAgent",
	}

	userRepo.On("FindByEmail", ctx, cmd.Email).Return(u, nil).Once()

	res, err := service.Login(ctx, cmd)
	require.Error(t, err)

	assert.Nil(t, res)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)

	issuer.AssertNotCalled(t, "Issue", mock.Anything)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestAuthService_Logout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	issuer := new(mocks.TokenIssuer)
	userRepo := new(mocks.UserRepository)
	sessionRepo := new(mocks.SessionRepository)
	pp := password.DefaultParams()
	refreshTokenTTL := 30 * 24 * time.Hour

	service := NewAuthService(issuer, userRepo, sessionRepo, pp, refreshTokenTTL)

	cmd := &ports.LogoutCmd{Token: "test-token"}
	hash := sha256.Sum256([]byte(cmd.Token))

	sessionRepo.On("Delete", ctx, domain.SessionToken(hash)).Return(nil).Once()

	err := service.Logout(ctx, cmd)
	require.NoError(t, err)

	sessionRepo.AssertExpectations(t)
}
