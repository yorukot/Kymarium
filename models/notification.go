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
	WebhookURL string `json:"webhook_url"`
}

// TelegramNotificationConfig describes the stored config for a Telegram notification channel.
type TelegramNotificationConfig struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}
