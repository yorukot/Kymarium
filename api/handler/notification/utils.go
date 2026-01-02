package notification

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yorukot/kymarium/models"
)

// ValidateNotificationConfig validates the config JSON for a notification type.
func validateNotificationConfig(notificationType models.NotificationType, raw json.RawMessage) error {
	if len(raw) == 0 {
		return fmt.Errorf("notification config is required")
	}

	v := validator.New()

	switch notificationType {
	case models.NotificationTypeDiscord:
		var cfg models.DiscordNotificationConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return fmt.Errorf("decode discord notification config: %w", err)
		}
		return v.Struct(cfg)
	case models.NotificationTypeSlack:
		var cfg models.SlackNotificationConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return fmt.Errorf("decode slack notification config: %w", err)
		}
		return v.Struct(cfg)
	case models.NotificationTypeTelegram:
		var cfg models.TelegramNotificationConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return fmt.Errorf("decode telegram notification config: %w", err)
		}
		return v.Struct(cfg)
	case models.NotificationTypeEmail:
		var cfg models.EmailNotificationConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return fmt.Errorf("decode email notification config: %w", err)
		}
		return v.Struct(cfg)
	default:
		return fmt.Errorf("unsupported notification type %q", notificationType)
	}
}
