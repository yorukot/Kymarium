package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/yorukot/kymarium/models"
)

// telegramAPIBase is overridable for testing.
var telegramAPIBase = "https://api.telegram.org"

func sendTelegram(ctx context.Context, client *http.Client, notification models.Notification, title, description string, _ models.PingStatus) error {
	var cfg models.TelegramNotificationConfig
	if err := json.Unmarshal(notification.Config, &cfg); err != nil {
		return fmt.Errorf("decode telegram config: %w", err)
	}

	if cfg.BotToken == "" || cfg.ChatID == "" {
		return errors.New("telegram bot_token and chat_id are required")
	}

	apiBase := strings.TrimSuffix(telegramAPIBase, "/")
	url := fmt.Sprintf("%s/bot%s/sendMessage", apiBase, cfg.BotToken)

	payload := map[string]interface{}{
		"chat_id": cfg.ChatID,
		"text":    strings.TrimSpace(fmt.Sprintf("%s\n\n%s", title, description)),
	}

	return postJSON(ctx, client, url, payload)
}
