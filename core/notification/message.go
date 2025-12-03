package notification

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yorukot/knocker/models"
)

// MessageInput captures the data used to build a notification message.
type MessageInput struct {
	MonitorName string
	Status      models.PingStatus
	Region      string
	LatencyMs   int
	CheckedAt   time.Time
	Detail      string
}

// FormatMessage generates a title and description for a notification.
func FormatMessage(input MessageInput) (string, string) {
	checkedAt := input.CheckedAt
	if checkedAt.IsZero() {
		checkedAt = time.Now().UTC()
	}

	title := fmt.Sprintf("%s is %s", input.MonitorName, strings.ToUpper(string(input.Status)))

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Monitor: %s\n", input.MonitorName))
	if input.Region != "" {
		builder.WriteString(fmt.Sprintf("Region: %s\n", input.Region))
	}
	builder.WriteString(fmt.Sprintf("Status: %s\n", strings.ToUpper(string(input.Status))))

	if input.LatencyMs > 0 {
		builder.WriteString(fmt.Sprintf("\nLatency: %dms", input.LatencyMs))
	}

	builder.WriteString(fmt.Sprintf("\nChecked at: %s", checkedAt.UTC().Format(time.RFC3339)))

	if detail := strings.TrimSpace(input.Detail); detail != "" {
		builder.WriteString(fmt.Sprintf("\n\nDetails: %s", detail))
	}

	return title, strings.TrimSpace(builder.String())
}

// DetailFromRaw extracts a human-readable detail string from the stored ping data.
func DetailFromRaw(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return strings.TrimSpace(string(raw))
	}

	if msg, ok := decoded["error"].(string); ok && msg != "" {
		return msg
	}

	if statusCode, ok := decoded["http_status_code"].(float64); ok {
		return fmt.Sprintf("HTTP status code: %d", int(statusCode))
	}

	pretty, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		return ""
	}

	return string(pretty)
}
