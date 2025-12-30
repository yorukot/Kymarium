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

// TeamMemberWithUser represents a team member along with user info.
type TeamMemberWithUser struct {
	TeamMember
	DisplayName string  `json:"display_name" db:"display_name"`
	Email       string  `json:"email" db:"email"`
	Avatar      *string `json:"avatar,omitempty" db:"avatar"`
}

type TeamInvite struct {
	ID           int64        `json:"id,string" db:"id"`
	TeamID       int64        `json:"team_id,string" db:"team_id"`
	InvitedBy    int64        `json:"invited_by,string" db:"invited_by"`
	InvitedTo    int64        `json:"invited_to,string" db:"invited_to"`
	InvitedEmail string       `json:"invited_email" db:"invited_email"`
	Role         MemberRole   `json:"role" db:"role"`
	Status       InviteStatus `json:"status" db:"status"`
	Token        *string      `json:"-" db:"token"`
	ExpiresAt    time.Time    `json:"expires_at" db:"expires_at"`
	AcceptedAt   *time.Time   `json:"accepted_at,omitempty" db:"accepted_at"`
	RejectedAt   *time.Time   `json:"rejected_at,omitempty" db:"rejected_at"`
	CanceledAt   *time.Time   `json:"canceled_at,omitempty" db:"canceled_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
}

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusRejected InviteStatus = "rejected"
	InviteStatusCanceled InviteStatus = "canceled"
)

// TeamWithRole represents a team along with the current member's role.
type TeamWithRole struct {
	Team
	Role MemberRole `json:"role" db:"role"`
}
