package password

import (
	"regexp"
	"testing"
)

func TestHash(t *testing.T) {
	t.Parallel()

	p := DefaultParams()
	password := "test-password"

	hash, err := Hash(password, p)
	if err != nil {
		t.Fatalf("unexpected error; %v", err)
	}

	wantFormat := regexp.MustCompile(
		`^\$argon2id\$v=\d+\$m=\d+,t=\d+,p=\d+\$[A-Za-z0-9+/]+\$[A-Za-z0-9+/]+$`,
	)

	if !wantFormat.MatchString(hash) {
		t.Errorf("want expre: %v,\nbat act hash: %v", wantFormat.String(), hash)
	}
}

func TestHash_DifferentSaltsProduceDifferentHashes(t *testing.T) {
	t.Parallel()

	p := DefaultParams()
	password := "test-password"

	h1, err := Hash(password, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h2, err := Hash(password, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h1 == h2 {
		t.Fatal("expected different hash, got same")
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	p := DefaultParams()
	password := "test_password"
	correct_hash, err := Hash(password, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testCases := []struct {
		name    string
		hash    string
		wantVal bool
		wantErr error
	}{
		{
			name:    "correct_hash",
			hash:    correct_hash,
			wantVal: true,
			wantErr: nil,
		},
		{
			name:    "empty",
			hash:    "",
			wantVal: false,
			wantErr: ErrInvalidHash,
		},
		{
			name:    "not_enough_parts",
			hash:    "$argon2id$v=19$m=7168,t=5,p=2$onlysalt",
			wantVal: false,
			wantErr: ErrInvalidHash,
		},
		{
			name:    "wrong_algorithm",
			hash:    "$argon2$v=19$m=7168,t=5,p=2$test$test",
			wantVal: false,
			wantErr: ErrInvalidHash,
		},
		{
			name:    "bat_version",
			hash:    "$argon2id$v=xx$m=7168,t=5,p=2$onlysalt",
			wantVal: false,
			wantErr: ErrInvalidHash,
		},
		{
			name:    "incompatible_version",
			hash:    "$argon2id$v=999$m=7168,t=5,p=2$salt$hash",
			wantVal: false,
			wantErr: ErrIncompatibleVersion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			act, err := Verify(password, tc.hash)

			if err != tc.wantErr {
				t.Errorf("wantErr: %v, actErr: %v", tc.wantErr, err)
			}

			if act != tc.wantVal {
				t.Errorf("want: %t, act: %t", tc.wantVal, act)
			}
		})
	}
}
