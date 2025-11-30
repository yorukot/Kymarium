package models

import "time"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id,string" db:"id" example:"175928847299117063"`                        // Unique identifier for the user
	PasswordHash *string   `json:"password_hash,omitempty" db:"password_hash" example:"hashed_password"`  // Hashed password (omitted in responses)
	DisplayName  string    `json:"display_name" db:"display_name" example:"John Doe"`                     // Display name for the user
	Avatar       *string   `json:"avatar,omitempty" db:"avatar" example:"https://example.com/avatar.jpg"` // URL to user's avatar image
	CreatedAt    time.Time `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"`             // Timestamp when the user was created
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" example:"2023-01-01T12:00:00Z"`             // Timestamp when the user was last updated
}

// Account represents how a user can login to the system
type Account struct {
	ID             int64     `json:"id,string" db:"id" example:"175928847299117063"`            // Unique identifier for the account
	Provider       Provider  `json:"provider" db:"provider" example:"email"`                    // Authentication provider type
	ProviderUserID string    `json:"provider_user_id" db:"provider_user_id" example:"user123"`  // User ID from the provider
	UserID         int64     `json:"user_id,string" db:"user_id" example:"175928847299117063"`  // Associated user ID
	Email          string    `json:"email" db:"email" example:"user@example.com"`               // User's email address
	CreatedAt      time.Time `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"` // Timestamp when the account was created
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at" example:"2023-01-01T12:00:00Z"` // Timestamp when the account was last updated
}

// OAuthToken represents OAuth tokens for external providers
type OAuthToken struct {
	AccountID    int64     `json:"account_id,string" db:"account_id" example:"175928847299117063"`         // Associated account ID
	AccessToken  string    `json:"access_token" db:"access_token" example:"ya29.a0AfH6SMC..."`             // OAuth access token
	RefreshToken *string   `json:"refresh_token" db:"refresh_token" example:"1//0GWthXqhYjIsKCgYIARAA..."` // OAuth refresh token
	Expiry       time.Time `json:"expiry" db:"expiry" example:"2023-01-01T13:00:00Z"`                      // Token expiration time
	TokenType    string    `json:"token_type" db:"token_type" example:"Bearer"`                            // Token type (usually Bearer)
	Provider     Provider  `json:"provider" db:"provider" example:"google"`                                // OAuth provider
	CreatedAt    time.Time `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"`              // Timestamp when the token was created
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at" example:"2023-01-01T12:00:00Z"`              // Timestamp when the token was last updated
}
