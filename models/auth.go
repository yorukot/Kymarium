package models

import (
	"net"
	"time"
)

// CookieName constants
const (
	CookieNameOAuthSession = "oauth_session"
	CookieNameRefreshToken = "refresh_token"
	CookieNameAccessToken  = "access_token"
)

// RefreshToken represents a persisted refresh token session.
type RefreshToken struct {
	ID        int64      `json:"id,string" db:"id" example:"175928847299117063"`
	UserID    int64      `json:"user_id,string" db:"user_id" example:"175928847299117063"`
	Token     string     `json:"token" db:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserAgent *string    `json:"user_agent" db:"user_agent" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`
	IP        net.IP     `json:"ip" db:"ip" example:"192.168.1.100"`
	UsedAt    *time.Time `json:"used_at,omitempty" db:"used_at" example:"2023-01-01T12:00:00Z"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"`
}

// Provider represents the authentication provider type
type Provider string

// Provider constants
const (
	ProviderEmail  Provider = "email"  // Email/password authentication
	ProviderGoogle Provider = "google" // Google OAuth authentication
	// You can add more providers here
)
