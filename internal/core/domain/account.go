package domain

import "time"

type Account struct {
	ID                    int
	UserId                int
	AccountId             string
	ProviderId            string
	AccessToken           *string
	RefreshToken          *string
	AccessTokenExpiresAt  *time.Time
	RefreshTokenExpiresAt *time.Time
	Scope                 *string
	IdToken               *string
	Password              *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
