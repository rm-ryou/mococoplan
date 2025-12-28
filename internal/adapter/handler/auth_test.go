package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rm-ryou/mococoplan/internal/core/ports"
	"github.com/rm-ryou/mococoplan/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_SuccessSignup(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	ttl := 30 * 24 * time.Hour
	as.On("RefreshTokenTTL").Return(ttl).Once()

	ua := "test-agent"
	requestName := "Test"
	requestEmail := "test@example.com"
	requestPassword := "test-password"

	cmd := &ports.SignupCmd{
		Name: requestName,
		LoginCmd: ports.LoginCmd{
			Email:     requestEmail,
			Password:  requestPassword,
			UserAgent: ua,
		},
	}

	expiresAt := time.Now().Add(15 * time.Minute).UTC()
	as.On("Signup", mock.Anything, cmd).Return(&ports.AuthResult{
		AccessToken:  "access.jwt",
		ExpiresAt:    expiresAt,
		RefreshToken: "refresh.plain",
	}, nil).Once()

	body := map[string]any{
		"name":     requestName,
		"email":    requestEmail,
		"password": requestPassword,
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", ua)

	rr := httptest.NewRecorder()
	h.Signup(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	assert.Contains(t, res.Header.Get("Content-Type"), "application/json")

	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "refresh_token" {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	assert.Equal(t, "refresh.plain", cookie.Value)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.Equal(t, "/", cookie.Path)

	now := time.Now()
	assert.True(t, cookie.Expires.After(now.Add(ttl-1*time.Minute)))
	assert.True(t, cookie.Expires.Before(now.Add(ttl+1*time.Minute)))

	var resBody map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&resBody))

	assert.Equal(t, "access.jwt", resBody["access_token"])
	assert.NotNil(t, resBody["expires_at"])

	as.AssertExpectations(t)
}

func TestAuthHandler_FailedSignup_InvalidJson(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Signup(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	as.AssertNotCalled(t, "Signup", mock.Anything, mock.Anything)
}

func TestAuthHandler_FailedSignup_ValidationError(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	testCases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "invalid email format",
			body: map[string]any{
				"name":     "Test",
				"email":    "invalid-email-format",
				"password": "test-password",
			},
		},
		{
			name: "name less than 50 char",
			body: map[string]any{
				"name":     strings.Repeat("a", 51),
				"email":    "test@example.com",
				"password": "test-password",
			},
		},
		{
			name: "password more than 12 char",
			body: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": "a",
			},
		},
		{
			name: "password less than 128 char",
			body: map[string]any{
				"name":     "Test",
				"email":    "test@example.com",
				"password": strings.Repeat("a", 129),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.Signup(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			as.AssertNotCalled(t, "Signup", mock.Anything, mock.Anything)
		})
	}
}

func TestAuthHandler_SuccessLogin(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	ttl := 30 * 24 * time.Hour
	as.On("RefreshTokenTTL").Return(ttl).Once()

	ua := "test-agent"
	requestEmail := "test@example.com"
	requestPassword := "test-password"

	cmd := &ports.LoginCmd{
		Email:     requestEmail,
		Password:  requestPassword,
		UserAgent: ua,
	}

	expiresAt := time.Now().Add(15 * time.Minute).UTC()
	as.On("Login", mock.Anything, cmd).Return(&ports.AuthResult{
		AccessToken:  "access.jwt",
		ExpiresAt:    expiresAt,
		RefreshToken: "refresh.plain",
	}, nil).Once()

	body := map[string]any{
		"email":    requestEmail,
		"password": requestPassword,
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", ua)

	rr := httptest.NewRecorder()
	h.Login(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "refresh_token" {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	assert.Equal(t, "refresh.plain", cookie.Value)
	assert.True(t, cookie.HttpOnly)
	assert.True(t, cookie.Secure)
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	assert.Equal(t, "/", cookie.Path)

	now := time.Now()
	assert.True(t, cookie.Expires.After(now.Add(ttl-1*time.Minute)))
	assert.True(t, cookie.Expires.Before(now.Add(ttl+1*time.Minute)))

	var resBody map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&resBody))

	assert.Equal(t, "access.jwt", resBody["access_token"])
	assert.NotNil(t, resBody["expires_at"])

	as.AssertExpectations(t)
}

func TestAuthHandler_FailedLogin_InvalidJson(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	as.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
}

func TestAuthHandler_FailedLogin_ValidationError(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	testCases := []struct {
		name string
		body map[string]any
	}{
		{
			name: "invalid email format",
			body: map[string]any{
				"email":    "invalid-email-format",
				"password": "test-password",
			},
		},
		{
			name: "password more than 12 char",
			body: map[string]any{
				"email":    "test@example.com",
				"password": "a",
			},
		},
		{
			name: "password less than 128 char",
			body: map[string]any{
				"email":    "test@example.com",
				"password": strings.Repeat("a", 129),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			h.Signup(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
			as.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
		})
	}
}

func TestAuthhandler_SuccessLogout(t *testing.T) {
	t.Parallel()

	as := new(mocks.AuthServicer)
	h := NewAuthHandler(as)

	as.On("Logout", mock.Anything, mock.Anything).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: "refresh.plain",
	})
	rr := httptest.NewRecorder()

	h.Logout(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	var cookie *http.Cookie
	for _, c := range res.Cookies() {
		if c.Name == "refresh_token" {
			cookie = c
			break
		}
	}
	require.NotNil(t, cookie)
	assert.Equal(t, "", cookie.Value)
	assert.Equal(t, -1, cookie.MaxAge)

	as.AssertExpectations(t)
}
