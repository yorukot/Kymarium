package handler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hibiken/asynq"
	notificationcore "github.com/yorukot/kymarium/core/notification"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/config"
	"github.com/yorukot/kymarium/worker/tasks"
	"go.uber.org/zap"
)

// HandleNotificationDispatch processes notification dispatch tasks.
func (h *Handler) HandleNotificationDispatch(ctx context.Context, t *asynq.Task) error {
	var payload tasks.NotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		zap.L().Error("invalid notification payload", zap.Error(err))
		return err
	}

	monitor, notification, err := h.fetchMonitorAndNotification(ctx, payload)
	if err != nil {
		zap.L().Error("failed to load notification context",
			zap.Int64("monitor_id", payload.MonitorID),
			zap.Int64("notification_id", payload.NotificationID),
			zap.Error(err))
		return err
	}

	if monitor == nil {
		zap.L().Warn("monitor not found for notification dispatch",
			zap.Int64("monitor_id", payload.MonitorID),
			zap.Int64("notification_id", payload.NotificationID))
		return nil
	}

	if notification == nil {
		zap.L().Warn("notification not found for dispatch",
			zap.Int64("monitor_id", payload.MonitorID),
			zap.Int64("notification_id", payload.NotificationID))
		return nil
	}

	detail := strings.TrimSpace(payload.Detail)
	region := config.RegionByID(payload.RegionID)
	title, description := notificationcore.FormatMessage(notificationcore.MessageInput{
		MonitorName:       monitor.Name,
		Status:            payload.Ping.Status,
		RegionDisplayName: region.DisplayName,
		LatencyMs:         payload.Ping.Latency,
		CheckedAt:         payload.Ping.Time,
		Detail:            detail,
	})
	if err := notificationcore.Send(ctx, *notification, title, description, payload.Ping.Status); err != nil {
		zap.L().Error("failed to send notification",
			zap.Int64("monitor_id", payload.MonitorID),
			zap.Int64("notification_id", payload.NotificationID),
			zap.String("notification_type", string(notification.Type)),
			zap.Error(err))
		return err
	}

	zap.L().Info("notification dispatched",
		zap.Int64("monitor_id", payload.MonitorID),
		zap.Int64("notification_id", payload.NotificationID),
		zap.String("notification_type", string(notification.Type)),
		zap.Int64("region_id", payload.RegionID),
		zap.String("region", region.Name),
		zap.String("status", string(payload.Ping.Status)))

	return nil
}

func (h *Handler) fetchMonitorAndNotification(ctx context.Context, payload tasks.NotificationPayload) (*models.Monitor, *models.Notification, error) {
	tx, err := h.repo.StartTransaction(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer h.repo.DeferRollback(ctx, tx)

	monitor, err := h.repo.GetMonitorByID(ctx, tx, payload.TeamID, payload.MonitorID)
	if err != nil || monitor == nil {
		return monitor, nil, err
	}

	notification, err := h.repo.GetNotificationByID(ctx, tx, payload.TeamID, payload.NotificationID)
	if err != nil || notification == nil {
		return monitor, notification, err
	}

	if err := h.repo.CommitTransaction(ctx, tx); err != nil {
		return nil, nil, err
	}

	return monitor, notification, nil
}
