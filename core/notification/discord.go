package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/yorukot/knocker/models"
)

func sendDiscord(ctx context.Context, client *http.Client, notification models.Notification, title, description string) error {
	var cfg models.DiscordNotificationConfig
	if err := json.Unmarshal(notification.Config, &cfg); err != nil {
		return fmt.Errorf("decode discord config: %w", err)
	}

	if cfg.WebhookURL == "" {
		return errors.New("discord webhook_url is required")
	}

	payload := map[string]any{
		"username": notification.Name,
		"embeds": []map[string]any{
			{
				"title":       title,
				"description": description,
			},
		},
	}

	return postJSON(ctx, client, cfg.WebhookURL, payload)
}
