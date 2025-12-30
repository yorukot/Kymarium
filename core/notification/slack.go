package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/yorukot/knocker/models"
)

func sendSlack(ctx context.Context, client *http.Client, notification models.Notification, title, description string, _ models.PingStatus) error {
	var cfg models.SlackNotificationConfig
	if err := json.Unmarshal(notification.Config, &cfg); err != nil {
		return fmt.Errorf("decode slack config: %w", err)
	}

	if cfg.WebhookURL == "" {
		return errors.New("slack webhook_url is required")
	}

	text := strings.TrimSpace(fmt.Sprintf("*%s*\n%s", title, description))
	payload := map[string]any{
		"text": text,
	}

	return postJSON(ctx, client, cfg.WebhookURL, payload)
}
