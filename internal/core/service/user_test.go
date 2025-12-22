package service

import (
	"context"
	"testing"

	"github.com/rm-ryou/mococoplan/internal/adapter/repository"
	"github.com/rm-ryou/mococoplan/internal/core/domain/user"
	"github.com/rm-ryou/mococoplan/internal/core/domain/user/mocks"
	"github.com/rm-ryou/mococoplan/pkg/password"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Create(t *testing.T) {
	mockRepo := new(mocks.Repository)
	params := &password.Params{
		Memory:      1,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  0,
		KeyLength:   1,
	}
	service := NewUserService(mockRepo, params)

	testCases := []struct {
		name    string
		mockCmd *user.CreateCmd
		setupFn func(t *testing.T, cmd *user.CreateCmd)
		wantErr error
	}{
		{
			name: "success",
			mockCmd: &user.CreateCmd{
				Name:          "test name",
				Email:         "test@example.com",
				PlainPassword: "testPassword",
			},
			setupFn: func(t *testing.T, cmd *user.CreateCmd) {
				t.Helper()

				hash, err := password.Hash(cmd.PlainPassword, params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				mockUser := &user.User{
					Name:         cmd.Name,
					Email:        cmd.Email,
					PasswordHash: hash,
				}
				mockRepo.On("Create", mock.Anything, mockUser).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "failed - when email already exists",
			mockCmd: &user.CreateCmd{
				Name:          "test name",
				Email:         "already-exists@example.com",
				PlainPassword: "testPassword",
			},
			setupFn: func(t *testing.T, cmd *user.CreateCmd) {
				t.Helper()

				hash, err := password.Hash(cmd.PlainPassword, params)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				mockUser := &user.User{
					Name:         cmd.Name,
					Email:        cmd.Email,
					PasswordHash: hash,
				}
				mockRepo.On("Create", mock.Anything, mockUser).Return(repository.ErrEmailAlreadyExists)
			},
			wantErr: repository.ErrEmailAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupFn(t, tc.mockCmd)

			err := service.Create(context.Background(), tc.mockCmd)
			if err != tc.wantErr {
				t.Errorf("want error %v, act: %v", tc.wantErr, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}

	t.Run("success", func(t *testing.T) {
		mockCmd := &user.CreateCmd{
			Name:          "test name",
			Email:         "test@example.com",
			PlainPassword: "testPassword",
		}

		hash, err := password.Hash(mockCmd.PlainPassword, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		mockUser := &user.User{
			Name:         mockCmd.Name,
			Email:        mockCmd.Email,
			PasswordHash: hash,
		}

		mockRepo.On("Create", mock.Anything, mockUser).Return(nil)

		err = service.Create(context.Background(), mockCmd)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		mockRepo.AssertExpectations(t)
	})

}
