package domain

import "time"

type Workspace struct {
	ID        int
	Name      string
	Slug      string
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkspaceRole string

const (
	WorkspaceRoleOwner  WorkspaceRole = "owner"
	WorkspaceRoleAdmin  WorkspaceRole = "admin"
	WorkspaceRoleMember WorkspaceRole = "member"
)

type WorkspaceMember struct {
	WorkspaceID int
	UserID      int
	Role        WorkspaceRole
	JoinedAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
