package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/config"
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
		subject = "Kymarium notification"
	}

	body := strings.TrimSpace(description)
	if body == "" {
		body = "No additional details provided."
	}

	return config.SendEmail(to, nil, bcc, subject, body)
}
