package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/config"
)

func sendEmail(ctx context.Context, _ *http.Client, notification models.Notification, title, description string, _ models.PingStatus) error {
	_ = ctx

	var cfg models.EmailNotificationConfig
	if err := json.Unmarshal(notification.Config, &cfg); err != nil {
		return fmt.Errorf("decode email config: %w", err)
	}

	to := cfg.EmailAddress[0]
	var bcc []string
	if len(cfg.EmailAddress) > 1 {
		bcc = cfg.EmailAddress[1:]
	}

	subject := strings.TrimSpace(title)
	if subject == "" {
		subject = "Knocker notification"
	}

	body := strings.TrimSpace(description)
	if body == "" {
		body = "No additional details provided."
	}

	return config.SendEmail(to, nil, bcc, subject, body)
}
