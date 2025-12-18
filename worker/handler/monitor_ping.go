package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	monitorcore "github.com/yorukot/knocker/core/monitor"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/utils/id"
	"github.com/yorukot/knocker/worker/tasks"
	"go.uber.org/zap"
)

// HandleStartServiceTask processes service start tasks.
func (h *Handler) HandleStartServiceTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.MonitorPingPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	ping, detail, err := h.pingMonitor(ctx, payload.Monitor, payload.RegionID)
	if err != nil {
		zap.L().Warn("monitor ping encountered error",
			zap.Int64("monitor_id", payload.Monitor.ID),
			zap.Int64("region_id", payload.RegionID),
			zap.Error(err))
	}

	h.pingBuffer.Record(ctx, ping)

	h.processIncident(ctx, payload.Monitor, ping, payload.RegionID, detail)

	// Errors are logged and captured in ping history; returning nil prevents repeated retries.
	return nil
}

func (h *Handler) pingMonitor(ctx context.Context, monitor models.Monitor, regionID int64) (models.Ping, string, error) {
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
		RegionID:  regionID,
		Status:    models.PingStatusFailed,
		Latency:   0,
	}

	if result != nil {
		ping.Status = result.Status
		ping.Latency = int(clampLatencyMs(result.Duration))
	}

	return ping, message, err
}

func (h *Handler) enqueueNotificationTasks(monitor models.Monitor, ping models.Ping, regionID int64, detail string) {
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
			RegionID:       regionID,
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

func (h *Handler) processIncident(ctx context.Context, monitor models.Monitor, ping models.Ping, regionID int64, detail string) {
	region := config.RegionByID(regionID)

	tx, err := h.repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("failed to start incident transaction",
			zap.Int64("monitor_id", monitor.ID),
			zap.Int64("region_id", regionID),
			zap.String("region_name", region.Name),
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

	// Update monitor status based on latest ping before incident logic.
	targetStatus := monitor.Status
	if ping.Status == models.PingStatusSuccessful {
		targetStatus = models.MonitorStatusUp
	} else {
		targetStatus = models.MonitorStatusDown
	}
	if targetStatus != monitor.Status {
		if err := h.repo.UpdateMonitorStatus(ctx, tx, monitor.ID, targetStatus, time.Now().UTC()); err != nil {
			zap.L().Error("failed to update monitor status",
			zap.Int64("monitor_id", monitor.ID),
			zap.Int64("region_id", regionID),
			zap.String("region_name", region.Name),
			zap.String("target_status", string(targetStatus)),
			zap.Error(err))
		return
	}
		monitor.Status = targetStatus
	}

	if ping.Status == models.PingStatusSuccessful {
		notify, notifyDetail, err = h.handleIncidentRecovery(ctx, tx, monitor, ping, regionID, detail, openIncident)
	} else {
		notify, notifyDetail, err = h.handleIncidentFailure(ctx, tx, monitor, ping, regionID, detail, openIncident)
	}

	if err != nil {
		zap.L().Error("incident handling failed",
			zap.Int64("monitor_id", monitor.ID),
			zap.Int64("region_id", regionID),
			zap.String("region_name", region.Name),
			zap.Error(err))
		return
	}

	if err := h.repo.CommitTransaction(tx, ctx); err != nil {
		zap.L().Error("failed to commit incident transaction",
			zap.Int64("monitor_id", monitor.ID),
			zap.Int64("region_id", regionID),
			zap.String("region_name", region.Name),
			zap.Error(err))
		return
	}

	if notify {
		h.enqueueNotificationTasks(monitor, ping, regionID, notifyDetail)
	}
}

func (h *Handler) handleIncidentFailure(ctx context.Context, tx pgx.Tx, monitor models.Monitor, ping models.Ping, regionID int64, detail string, openIncident *models.Incident) (bool, string, error) {
	// Maintain only one active incident per monitor; use the region-specific window for detection.
	failureThreshold := int(monitor.FailureThreshold)
	if failureThreshold <= 0 {
		return false, "", nil
	}

	window := int(math.Ceil(float64(failureThreshold) * 1.5))
	recent, err := h.repo.ListRecentPingsByMonitorIDAndRegion(ctx, tx, monitor.ID, regionID, window-1)
	if err != nil {
		return false, "", err
	}

	samples := append([]models.Ping{ping}, recent...)
	failureCount := countFailures(samples, window)

	now := time.Now().UTC()
	message := incidentMessage(strconv.FormatInt(regionID, 10), detail, ping, string(ping.Status))

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
		lastEvent, err := h.repo.GetLastEventTimeline(ctx, tx, openIncident.ID)
		if err != nil {
			return false, "", err
		}

		if lastEvent == nil || strings.TrimSpace(lastEvent.Message) != message {
			if err := h.repo.CreateEventTimeline(ctx, tx, models.EventTimeline{
				IncidentID: openIncident.ID,
				Message:    message,
				EventType:  models.IncidentEventTypeUpdate,
				CreatedAt:  now,
				UpdatedAt:  now,
			}); err != nil {
				return false, "", err
			}
		}
	}

	return false, "", nil
}

func (h *Handler) handleIncidentRecovery(ctx context.Context, tx pgx.Tx, monitor models.Monitor, ping models.Ping, regionID int64, detail string, openIncident *models.Incident) (bool, string, error) {
	// Nothing to do if no incident is open.
	if openIncident == nil {
		return false, "", nil
	}

	recoveryThreshold := int(monitor.RecoveryThreshold)
	if recoveryThreshold <= 0 {
		return false, "", nil
	}

	// Pull only enough recent pings (region-specific) to evaluate recovery, include current ping first.
	recent, err := h.repo.ListRecentPingsByMonitorIDAndRegion(ctx, tx, monitor.ID, regionID, recoveryThreshold-1)
	if err != nil {
		return false, "", err
	}

	samples := append([]models.Ping{ping}, recent...)
	if len(samples) < recoveryThreshold {
		return false, "", nil
	}

	allSuccessful := true
	for i := 0; i < recoveryThreshold; i++ {
		if samples[i].Status != models.PingStatusSuccessful {
			allSuccessful = false
			break
		}
	}

	if !allSuccessful {
		return false, "", nil
	}

	now := time.Now().UTC()
	message := incidentMessage(strconv.FormatInt(regionID, 10), detail, ping, "recovered")

	if err := h.repo.MarkIncidentResolved(ctx, tx, openIncident.ID, ping.Time, now); err != nil {
		return false, "", err
	}

	if err := h.repo.CreateEventTimeline(ctx, tx, models.EventTimeline{
		IncidentID: openIncident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeAutoResolved,
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
	for i := range limit {
		if pings[i].Status != models.PingStatusSuccessful {
			failures++
		}
	}

	return failures
}

// createIncidentIfAbsent tries to create a new incident, handling unique constraint races gracefully.
func (h *Handler) createIncidentIfAbsent(ctx context.Context, tx pgx.Tx, monitorID int64, startedAt time.Time, message string, now time.Time) (*models.Incident, bool, error) {
	newID, err := id.GetID()
	if err != nil {
		return nil, false, err
	}

	incident := models.Incident{
		ID:          newID,
		Status:      models.IncidentStatusDetected,
		Severity:    models.IncidentSeverityMajor,
		IsPublic:    false,
		AutoResolve: true,
		StartedAt:   startedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
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

	if err := h.repo.CreateIncidentMonitor(ctx, tx, incident.ID, monitorID); err != nil {
		return nil, false, err
	}

	if err := h.repo.CreateEventTimeline(ctx, tx, models.EventTimeline{
		IncidentID: incident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeDetected,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, false, err
	}

	if err := h.repo.CreateEventTimeline(ctx, tx, models.EventTimeline{
		IncidentID: incident.ID,
		Message:    message,
		EventType:  models.IncidentEventTypeNotificationSent,
		CreatedAt:  now,
		UpdatedAt:  now,
	}); err != nil {
		return nil, false, err
	}

	return &incident, true, nil
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
