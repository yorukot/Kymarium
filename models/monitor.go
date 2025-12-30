package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/yorukot/knocker/models/monitorm"
)

type MonitorType string

const (
	MonitorTypeHTTP MonitorType = "http"
	MonitorTypePing MonitorType = "ping"
)

type MonitorStatus string

const (
	MonitorStatusUp   MonitorStatus = "up"
	MonitorStatusDown MonitorStatus = "down"
)

type NotificationType string

const (
	NotificationTypeDiscord  NotificationType = "discord"
	NotificationTypeTelegram NotificationType = "telegram"
	NotificationTypeSlack    NotificationType = "slack"
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
	Status   MonitorStatus   `json:"status" db:"status"`

	// Scheduling
	LastChecked time.Time `json:"last_checked" db:"last_checked"`
	NextCheck   time.Time `json:"next_check" db:"next_check"`

	// Thresholds
	FailureThreshold  int16 `json:"failure_threshold" db:"failure_threshold"`
	RecoveryThreshold int16 `json:"recovery_threshold" db:"recovery_threshold"`

	// Regions
	RegionIDs []int64 `json:"regions" db:"region_ids"`

	// Notifications
	NotificationIDs []int64 `json:"notification" db:"notification_ids"`

	// Metadata
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// MonitorWithIncidents decorates a monitor with its recent incidents for list views.
type MonitorWithIncidents struct {
	Monitor
	Incidents []Incident `json:"incidents,omitempty" db:"incidents"`
}

// HTTPConfig decodes the monitor config into an HTTPMonitorConfig.
func (m Monitor) HTTPConfig() (*monitorm.HTTPMonitorConfig, error) {
	if m.Type != MonitorTypeHTTP {
		return nil, fmt.Errorf("unsupported monitor type %q", m.Type)
	}

	var cfg monitorm.HTTPMonitorConfig
	if err := json.Unmarshal(m.Config, &cfg); err != nil {
		return nil, fmt.Errorf("decode http monitor config: %w", err)
	}

	return &cfg, validator.New().Struct(cfg)
}

// PingConfig decodes the monitor config into a PingMonitorConfig.
func (m Monitor) PingConfig() (*monitorm.PingMonitorConfig, error) {
	if m.Type != MonitorTypePing {
		return nil, fmt.Errorf("unsupported monitor type %q", m.Type)
	}

	var cfg monitorm.PingMonitorConfig
	if err := json.Unmarshal(m.Config, &cfg); err != nil {
		return nil, fmt.Errorf("decode ping monitor config: %w", err)
	}

	return &cfg, validator.New().Struct(cfg)
}
