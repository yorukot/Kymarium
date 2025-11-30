package models

import "time"

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
	MemberRoleViewer MemberRole = "viewer"
)

type Team struct {
	ID        int64     `json:"id,string" db:"id"`
	Name      string    `json:"name" db:"name"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TeamMember struct {
	ID        int64      `json:"id,string" db:"id"`
	TeamID    int64      `json:"team_id,string" db:"team_id"`
	UserID    int64      `json:"user_id,string" db:"user_id"`
	Role      MemberRole `json:"role" db:"role"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type TeamInvite struct {
	ID        int64     `json:"id,string" db:"id"`
	TeamID    int64     `json:"team_id,string" db:"team_id"`
	InvitedBy int64     `json:"invited_by,string" db:"Invited_by"`
	InvitedTo int64     `json:"invited_to,string" db:"Invited_to"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TeamWithRole represents a team along with the current member's role.
type TeamWithRole struct {
	Team
	Role MemberRole `json:"role" db:"role"`
}
