package domain

import "time"

type Account struct {
	ID           int
	UserId       int
	AccountId    string
	ProviderId   string
	AccessToken  *string
	RefreshToken *string
	IdToken      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
