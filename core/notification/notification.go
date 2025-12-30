package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yorukot/knocker/models"
)

// Send dispatches a notification using the provided notification model.
// Title and description are sent to the configured channel depending on the notification type.
func Send(ctx context.Context, notification models.Notification, title, description string, status models.PingStatus) error {
	return SendWithClient(ctx, http.DefaultClient, notification, title, description, status)
}

// SendWithClient allows injecting a custom HTTP client (useful for tests) while sending the notification.
func SendWithClient(ctx context.Context, client *http.Client, notification models.Notification, title, description string, status models.PingStatus) error {
	if client == nil {
		client = http.DefaultClient
	}

	switch notification.Type {
	case models.NotificationTypeDiscord:
		return sendDiscord(ctx, client, notification, title, description, status)
	case models.NotificationTypeSlack:
		return sendSlack(ctx, client, notification, title, description, status)
	case models.NotificationTypeTelegram:
		return sendTelegram(ctx, client, notification, title, description, status)
	case models.NotificationTypeEmail:
		return sendEmail(ctx, client, notification, title, description, status)
	default:
		return fmt.Errorf("unsupported notification type %q", notification.Type)
	}
}

func postJSON(ctx context.Context, client *http.Client, url string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return fmt.Errorf("unexpected status %d from %s: %s", resp.StatusCode, url, strings.TrimSpace(string(respBody)))
}
