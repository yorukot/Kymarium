package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type MonitorType string

const (
	MonitorTypeHTTP MonitorType = "http"
)

type NotificationType string

const (
	NotificationTypeDiscord  NotificationType = "discord"
	NotificationTypeTelegram NotificationType = "telegram"
	NotificationTypeEmail    NotificationType = "email"
)

// Monitor represents a monitor entity in the database.
// Fields are ordered by importance: identity, configuration, scheduling, notifications, metadata.
type Monitor struct {
	// Identity fields
	ID     int64 `json:"id,string" db:"id"`
	TeamID int64 `json:"team_id,string" db:"team_id"`

	// Core configuration
	Name     string          `json:"name" db:"name"`
	Type     MonitorType     `json:"type" db:"type"`
	Config   json.RawMessage `json:"config" db:"config"`
	Interval int             `json:"interval" db:"interval"`

	// Scheduling
	LastChecked time.Time `json:"last_checked" db:"last_checked"`
	NextCheck   time.Time `json:"next_check" db:"next_check"`

	// Notifications
	NotificationIDs []int64 `json:"notification" db:"notification"`

	// Metadata
	GroupID   *int64    `json:"group,omitempty" db:"group"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type HTTPMethod string

const (
    MethodGet     HTTPMethod = http.MethodGet
    MethodPost    HTTPMethod = http.MethodPost
    MethodPut     HTTPMethod = http.MethodPut
    MethodDelete  HTTPMethod = http.MethodDelete
    MethodPatch   HTTPMethod = http.MethodPatch
    MethodHead    HTTPMethod = http.MethodHead
    MethodOptions HTTPMethod = http.MethodOptions
)

// HTTPMonitorConfig represents the expected config shape for HTTP monitors.
// Fields are ordered by importance and functional grouping.
type HTTPMonitorConfig struct {
	// Core request configuration
	URL       string `json:"url"`
	Method    HTTPMethod `json:"method"`
	MaxRedirs int    `json:"max_redirects"`

	// Request options
	RequestTimeout int               `json:"request_timeout"`
	Headers        map[string]string `json:"headers,omitempty"`
	BodyEncoding   string            `json:"body_encoding,omitempty"`
	Body           string            `json:"body,omitempty"`

	// Response validation
	UpSideDownMode                bool   `json:"upside_down_mode"`
	CertificateExpiryNotification bool   `json:"certificate_expiry_notification"`
	IgnoreTLSError                bool   `json:"ignore_tls_error"`
	AcceptedStatusCodes           []int  `json:"accepted_status_codes"`

	// Notification options
	ResendThreshold int `json:"resend_threshold"`
}

// HTTPConfig decodes the monitor config into an HTTPMonitorConfig.
func (m Monitor) HTTPConfig() (*HTTPMonitorConfig, error) {
	if m.Type != MonitorTypeHTTP {
		return nil, fmt.Errorf("unsupported monitor type %q", m.Type)
	}

	var cfg HTTPMonitorConfig
	if err := json.Unmarshal(m.Config, &cfg); err != nil {
		return nil, fmt.Errorf("decode http monitor config: %w", err)
	}

	if cfg.URL == "" {
		return nil, errors.New("http monitor config missing url")
	}

	return &cfg, nil
}
