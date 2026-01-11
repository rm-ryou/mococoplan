package domain

import "time"

type Channel struct {
	ID          int
	WorkspaceID int
	Name        string
	Description string
	IsPrivate   bool
	CreatedBy   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ChannelWithJoin struct {
	Channel
	Joined bool
}

type ChannelMember struct {
	ChannelID int
	UserID    int
	JoinedAt  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
