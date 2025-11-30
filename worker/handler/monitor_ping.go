package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/models"
	"go.uber.org/zap"
)

// HandleStartServiceTask processes service start tasks.
func (h *Handler) HandleStartServiceTask(ctx context.Context, t *asynq.Task) error {
	var monitor models.Monitor
	if err := json.Unmarshal(t.Payload(), &monitor); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	cfg, err := monitor.HTTPConfig()
	if err != nil {
		return fmt.Errorf("monitor config invalid: %w", err)
	}

	// Use get to the monitor URL
	resp, err := http.Get(cfg.URL)
	if err != nil {
		zap.L().Error("Error making GET request", zap.Error(err))
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	zap.L().Debug("Ping completed",
		zap.Int64("monitor_id", monitor.ID),
		zap.String("url", cfg.URL),
		zap.Int("status_code", resp.StatusCode),
		zap.String("status", resp.Status))

	// Start service
	return nil
}
