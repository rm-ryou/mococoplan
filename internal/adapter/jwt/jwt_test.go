package jwt

import (
	"testing"
	"time"

	"github.com/rm-ryou/mococoplan/internal/core/token"
)

func Test_IssueAndVerify(t *testing.T) {
	t.Parallel()

	svc := New("secret", "testIsuuer", 1*time.Minute)

	in := &token.Claims{
		UserId: 1,
		Email:  "test@example.com",
	}

	at, err := svc.Issue(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if at.Token == "" {
		t.Fatal("token is not empty")
	}
	if at.ExpiresAt.Before(time.Now()) {
		t.Fatalf("expiresAt should be future: %v", at.ExpiresAt)
	}

	out, err := svc.Verify(at.Token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.UserId != in.UserId || out.Email != in.Email {
		t.Fatalf("want: %v, act: %v", out, in)
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	svc := New("secret", "test", 1*time.Minute)

	tc := &token.Claims{
		UserId: 1,
		Email:  "test@example.com",
	}

	at, err := svc.Issue(tc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testCases := []struct {
		name    string
		token   string
		wantVal *token.Claims
		wantErr error
	}{
		{
			name:    "invalid token",
			token:   "this.token.is.not.jwt",
			wantVal: nil,
			wantErr: ErrInvalidToken,
		},
		{
			name:    "tampered token",
			token:   at.Token + "tampered",
			wantVal: nil,
			wantErr: ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act, err := svc.Verify(tc.token)

			if act != tc.wantVal {
				t.Errorf("want: %v, act: %v", tc.wantVal, act)
			}

			if err != tc.wantErr {
				t.Errorf("wantErr: %v, act: %v", tc.wantErr, err)
			}
		})
	}
}
