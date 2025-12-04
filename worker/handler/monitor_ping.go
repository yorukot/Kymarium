package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	ping, detail, err := h.pingMonitor(ctx, payload.Monitor, payload.Region)
	if err != nil {
		zap.L().Warn("monitor ping encountered error",
			zap.Int64("monitor_id", payload.Monitor.ID),
			zap.String("region", payload.Region),
			zap.Error(err))
	}

	h.pingBuffer.Record(ctx, ping)

	h.processIncident(ctx, payload.Monitor, ping, payload.Region, detail)

	// Errors are logged and captured in ping history; returning nil prevents repeated retries.
	return nil
}

func (h *Handler) pingMonitor(ctx context.Context, monitor models.Monitor, region string) (models.Ping, string, error) {
	result, err := monitorcore.Run(ctx, monitor)

	message := ""
	if result != nil {
		message = result.Message
	} else if err != nil {
		message = err.Error()
	}

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
	}

	return ping, message, err
}

func (h *Handler) enqueueNotificationTasks(monitor models.Monitor, ping models.Ping, region string, detail string) {
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
			Detail:         detail,
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

func (h *Handler) processIncident(ctx context.Context, monitor models.Monitor, ping models.Ping, region, detail string) {
	tx, err := h.repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("failed to start incident transaction",
			zap.Int64("monitor_id", monitor.ID),
			zap.String("region", region),
			zap.Error(err))
		return
	}
	defer h.repo.DeferRollback(tx, ctx)

	openIncident, err := h.repo.GetOpenIncidentByMonitorID(ctx, tx, monitor.ID)
	if err != nil {
		zap.L().Error("failed to load open incident",
			zap.Int64("monitor_id", monitor.ID),
			zap.Error(err))
		return
	}

	var notify bool
	var notifyDetail string

	if ping.Status == models.PingStatusSuccessful {
		notify, notifyDetail, err = h.handleIncidentRecovery(ctx, tx, monitor, ping, region, detail, openIncident)
	} else {
		notify, notifyDetail, err = h.handleIncidentFailure(ctx, tx, monitor, ping, region, detail, openIncident)
	}

	if err != nil {
		zap.L().Error("incident handling failed",
			zap.Int64("monitor_id", monitor.ID),
			zap.String("region", region),
			zap.Error(err))
		return
	}

	if err := h.repo.CommitTransaction(tx, ctx); err != nil {
		zap.L().Error("failed to commit incident transaction",
			zap.Int64("monitor_id", monitor.ID),
			zap.String("region", region),
			zap.Error(err))
		return
	}

	if notify {
		h.enqueueNotificationTasks(monitor, ping, region, notifyDetail)
	}
}

func (h *Handler) handleIncidentFailure(ctx context.Context, tx pgx.Tx, monitor models.Monitor, ping models.Ping, region, detail string, openIncident *models.Incident) (bool, string, error) {
	// Maintain only one active incident per monitor; use the region-specific window for detection.
	failureThreshold := int(monitor.FailureThreshold)
	if failureThreshold <= 0 {
		return false, "", nil
	}

	window := int(math.Ceil(float64(failureThreshold) * 1.5))
	recent, err := h.repo.ListRecentPingsByMonitorIDAndRegion(ctx, tx, monitor.ID, region, window-1)
	if err != nil {
		return false, "", err
	}

	samples := append([]models.Ping{ping}, recent...)
	failureCount := countFailures(samples, window)

	now := time.Now().UTC()
	message := incidentMessage(region, detail, ping, string(ping.Status))

	// Create a new incident when the failure threshold is met.
	if failureCount >= failureThreshold && len(samples) >= failureThreshold && openIncident == nil {
		createdIncident, created, err := h.createIncidentIfAbsent(ctx, tx, monitor.ID, ping.Time, message, now)
		if err != nil {
			return false, "", err
		}
		if created {
			return true, message, nil
		}
		// If not created, fall through to update handling below.
		openIncident = createdIncident
	}

	// If an incident is already open, record meaningful changes in failure type.
	if openIncident != nil {
		lastEvent, err := h.repo.GetLastIncidentEvent(ctx, tx, openIncident.ID)
		if err != nil {
			return false, "", err
		}

		if lastEvent == nil || strings.TrimSpace(lastEvent.Message) != message {
			if err := h.repo.CreateIncidentEvent(ctx, tx, models.IncidentEvent{
				IncidentID: openIncident.ID,
				Message:    message,
				EventType:  models.IncidentEventTypeUpdate,
				Public:     true,
				CreatedAt:  now,
				UpdatedAt:  now,
			}); err != nil {
				return false, "", err
			}
		}
	}

	return false, "", nil
}

func (h *Handler) handleIncidentRecovery(ctx context.Context, tx pgx.Tx, monitor models.Monitor, ping models.Ping, region, detail string, openIncident *models.Incident) (bool, string, error) {
	// Nothing to do if no incident is open.
	if openIncident == nil {
		return false, "", nil
	}

	recoveryThreshold := int(monitor.RecoveryThreshold)
	if recoveryThreshold <= 0 {
		return false, "", nil
	}

	// Pull only enough recent pings (region-specific) to evaluate recovery, include current ping first.
	recent, err := h.repo.ListRecentPingsByMonitorIDAndRegion(ctx, tx, monitor.ID, region, recoveryThreshold-1)
	if err != nil {
		return false, "", err
	}

	samples := append([]models.Ping{ping}, recent...)
	if len(samples) < recoveryThreshold {
		return false, "", nil
	}

	allSuccessful := true
	for i := range(recoveryThreshold) {
		if samples[i].Status != models.PingStatusSuccessful {
			allSuccessful = false
			break
		}
	}

	if !allSuccessful {
		return false, "", nil
	}

	now := time.Now().UTC()
	message := incidentMessage(region, detail, ping, "recovered")

	if err := h.repo.MarkIncidentResolved(ctx, tx, openIncident.ID, ping.Time, now); err != nil {
		return false, "", err
	}

	if err := h.repo.CreateIncidentEvent(ctx, tx, models.IncidentEvent{
		IncidentID: openIncident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeAutoResolved,
		Public:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return false, "", err
	}

	return true, message, nil
}

func countFailures(pings []models.Ping, window int) int {
	if window <= 0 {
		return 0
	}

	limit := min(len(pings), window)

	failures := 0
	for i := range(limit) {
		if pings[i].Status != models.PingStatusSuccessful {
			failures++
		}
	}

	return failures
}

// createIncidentIfAbsent tries to create a new incident, handling unique constraint races gracefully.
func (h *Handler) createIncidentIfAbsent(ctx context.Context, tx pgx.Tx, monitorID int64, startedAt time.Time, message string, now time.Time) (*models.Incident, bool, error) {
	incident := models.Incident{
		MonitorID: monitorID,
		Status:    models.IncidentStatusDetected,
		StartedAt: startedAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.repo.CreateIncident(ctx, tx, incident); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			existing, getErr := h.repo.GetOpenIncidentByMonitorID(ctx, tx, monitorID)
			if getErr != nil {
				return nil, false, getErr
			}
			if existing == nil {
				return nil, false, fmt.Errorf("unique violation but no open incident found")
			}
			return existing, false, nil
		}
		return nil, false, err
	}

	createdIncident, err := h.repo.GetOpenIncidentByMonitorID(ctx, tx, monitorID)
	if err != nil {
		return nil, false, err
	}
	if createdIncident == nil {
		return nil, false, fmt.Errorf("incident created but not found")
	}

	if err := h.repo.CreateIncidentEvent(ctx, tx, models.IncidentEvent{
		IncidentID: createdIncident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeDetected,
		Public:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, false, err
	}

	if err := h.repo.CreateIncidentEvent(ctx, tx, models.IncidentEvent{
		IncidentID: createdIncident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeNotificationSent,
		Public:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, false, err
	}

	return createdIncident, true, nil
}

func incidentMessage(region, detail string, ping models.Ping, fallback string) string {
	msg := strings.TrimSpace(detail)
	if msg == "" {
		msg = strings.TrimSpace(fallback)
	}
	if msg == "" {
		msg = fmt.Sprintf("status %s", ping.Status)
	}
	if region != "" {
		msg = fmt.Sprintf("%s: %s", region, msg)
	}
	return msg
}
