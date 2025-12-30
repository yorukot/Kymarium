package models

import (
	"encoding/json"
	"time"
)

type Notification struct {
	ID        int64            `json:"id,string" db:"id"`
	TeamID    int64            `json:"team_id,string" db:"team_id"`
	Type      NotificationType `json:"type" db:"type"`
	Name      string           `json:"name" db:"name"`
	Config    json.RawMessage  `json:"config" db:"config"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

// DiscordNotificationConfig describes the stored config for a Discord notification channel.
type DiscordNotificationConfig struct {
	WebhookURL string `json:"webhook_url" validate:"required,url"`
}

// SlackNotificationConfig describes the stored config for a Slack notification channel.
type SlackNotificationConfig struct {
	WebhookURL string `json:"webhook_url" validate:"required,url"`
}

// TelegramNotificationConfig describes the stored config for a Telegram notification channel.
type TelegramNotificationConfig struct {
	BotToken string `json:"bot_token" validate:"required,max=500"`
	ChatID   string `json:"chat_id" validate:"required,max=100"`
}

type EmailNotificationConfig struct {
	EmailAddress []string `json:"email_address" validate:"required,min=1,dive,required,email"`
}

type MonitorNotification struct {
	ID             int64 `json:"id,string" db:"id"`
	MonitorID      int64 `json:"monitor_id,string" db:"monitor_id"`
	NotificationID int64 `json:"notification_id,string" db:"notification_id"`
}
