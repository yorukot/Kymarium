package models

import (
	"encoding/json"
	"errors"
	"fmt"
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

// Monitor represents a monitor entity in the database
type Monitor struct {
	ID              int64           `json:"id,string" db:"id"`
	TeamID          int64           `json:"team_id,string" db:"team_id"`
	Name            string          `json:"name" db:"name"`
	Type            MonitorType     `json:"type" db:"type"`
	Interval        int             `json:"interval" db:"interval"`
	Config          json.RawMessage `json:"config" db:"config"`
	LastChecked     time.Time       `json:"last_checked" db:"last_checked"`
	NextCheck       time.Time       `json:"next_check" db:"next_check"`
	NotificationIDs []int64         `json:"notification" db:"notification"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
	CreatedAt       time.Time       `json:"created_at" db:"creted_at"`
	GroupID         *int64          `json:"group,omitempty" db:"group"`
}

// HTTPMonitorConfig represents the expected config shape for HTTP monitors.
type HTTPMonitorConfig struct {
	URL    string `json:"url"`
	Method string `json:"method,omitempty"`
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
