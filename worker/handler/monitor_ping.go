package handler

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/hibiken/asynq"
	monitorcore "github.com/yorukot/knocker/core/monitor"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

// HandleStartServiceTask processes service start tasks.
func (h *Handler) HandleStartServiceTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.MonitorPingPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	ping, err := h.pingMonitor(ctx, payload.Monitor, payload.Region)
	if err != nil {
		zap.L().Warn("monitor ping encountered error",
			zap.Int64("monitor_id", payload.Monitor.ID),
			zap.String("region", payload.Region),
			zap.Error(err))
	}

	h.pingBuffer.Record(ctx, ping)

	// Before enqueueing notifications, check if the ping was unsuccessful. and if this is the first time

	if ping.Status != models.PingStatusSuccessful {
		h.enqueueNotificationTasks(payload.Monitor, ping, payload.Region)
	}

	// Errors are logged and captured in ping history; returning nil prevents repeated retries.
	return nil
}

func (h *Handler) pingMonitor(ctx context.Context, monitor models.Monitor, region string) (models.Ping, error) {
	result, err := monitorcore.Run(ctx, monitor)

	ping := models.Ping{
		Time:      time.Now().UTC(),
		MonitorID: monitor.ID,
		Region:    region,
		Status:    models.PingStatusFailed,
		Latency:   0,
	}

	if result != nil {
		ping.Status = result.Status
		ping.Latency = int(clampLatencyMs(result.Duration))
		ping.Data = marshalPingData(result.Data)
	} else {
		ping.Data = marshalPingData(map[string]any{"error": err.Error()})
	}

	return ping, err
}

func (h *Handler) enqueueNotificationTasks(monitor models.Monitor, ping models.Ping, region string) {
	if h.notifier == nil {
		return
	}

	// Fetch notification IDs from junction table
	ctx := context.Background()
	tx, err := h.repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("failed to start transaction for notification fetch",
			zap.Int64("monitor_id", monitor.ID),
			zap.Error(err))
		return
	}
	defer h.repo.DeferRollback(tx, ctx)

	notificationIDs, err := h.repo.GetNotificationIDsByMonitorID(ctx, tx, monitor.ID)
	if err != nil {
		zap.L().Error("failed to fetch notification IDs",
			zap.Int64("monitor_id", monitor.ID),
			zap.Error(err))
		return
	}

	if err := h.repo.CommitTransaction(tx, ctx); err != nil {
		zap.L().Error("failed to commit transaction",
			zap.Int64("monitor_id", monitor.ID),
			zap.Error(err))
		return
	}

	if len(notificationIDs) == 0 {
		return
	}

	for _, notificationID := range notificationIDs {

		payload := tasks.NotificationPayload{
			TeamID:         monitor.TeamID,
			MonitorID:      monitor.ID,
			NotificationID: notificationID,
			Region:         region,
			Ping:           ping,
		}

		task, err := tasks.NewNotificationDispatch(payload)
		if err != nil {
			zap.L().Error("failed to create notification task",
				zap.Int64("monitor_id", monitor.ID),
				zap.Int64("notification_id", notificationID),
				zap.Error(err))
			continue
		}

		if _, err := h.notifier.Enqueue(task); err != nil {
			zap.L().Error("failed to enqueue notification task",
				zap.Int64("monitor_id", monitor.ID),
				zap.Int64("notification_id", notificationID),
				zap.Error(err))
		}
	}
}

func marshalPingData(data any) json.RawMessage {
	encoded, _ := json.Marshal(data)
	return encoded
}

func clampLatencyMs(duration time.Duration) int64 {
	ms := duration.Milliseconds()
	if ms < 0 {
		return 0
	}

	const maxLatencyMs = int64(math.MaxInt32)
	if ms > maxLatencyMs {
		return maxLatencyMs
	}
	return ms
}
