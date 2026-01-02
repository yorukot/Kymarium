package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/yorukot/kymarium/models"
)

func sendDiscord(ctx context.Context, client *http.Client, notification models.Notification, title, description string, status models.PingStatus) error {
	var cfg models.DiscordNotificationConfig
	if err := json.Unmarshal(notification.Config, &cfg); err != nil {
		return fmt.Errorf("decode discord config: %w", err)
	}

	if cfg.WebhookURL == "" {
		return errors.New("discord webhook_url is required")
	}

	// Discord doesn't allow usernames containing "discord" (case-insensitive)
	username := sanitizeDiscordUsername(notification.Name)

	payload := map[string]any{
		"username": username,
		"embeds": []map[string]any{
			{
				"title":       title,
				"description": description,
				"color":       discordColorForStatus(status),
			},
		},
	}

	return postJSON(ctx, client, cfg.WebhookURL, payload)
}

func sanitizeDiscordUsername(name string) string {
	re := regexp.MustCompile(`(?i)discord`)
	clean := strings.TrimSpace(re.ReplaceAllString(name, ""))
	if clean == "" {
		return "Kymarium"
	}
	return clean
}

func discordColorForStatus(status models.PingStatus) int {
	switch status {
	case models.PingStatusSuccessful:
		return 0x2ecc71 // green
	case models.PingStatusTimeout:
		return 0xf1c40f // yellow
	default:
		return 0xe74c3c // red
	}
}
