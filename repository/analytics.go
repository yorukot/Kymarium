package repository

import (
	"context"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// GetMonitorAnalytics retrieves aggregated uptime/latency buckets for a monitor over a time window.
// Data comes from the Timescale continuous aggregate monitor_30min_summary.
func (r *PGRepository) GetMonitorAnalytics(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time, regionID *int64) ([]models.MonitorAnalyticsBucket, error) {
	query := strings.Builder{}
	query.WriteString(`
		SELECT
			bucket,
			region_id,
			total_count,
			good_count,
			p50_ms,
			p75_ms,
			p90_ms,
			p95_ms,
			p99_ms
		FROM monitor_30min_summary
		WHERE monitor_id = $1
		  AND bucket >= $2
		  AND bucket < $3
	`)

	args := []any{monitorID, start, end}
	if regionID != nil {
		query.WriteString(" AND region_id = $4")
		args = append(args, *regionID)
	}

	query.WriteString(" ORDER BY bucket, region_id")

	var buckets []models.MonitorAnalyticsBucket
	if err := pgxscan.Select(ctx, tx, &buckets, query.String(), args...); err != nil {
		return nil, err
	}

	return buckets, nil
}

// ListIncidentsByMonitorIDWithinRange returns incidents overlapping the provided window for a monitor.
func (r *PGRepository) ListIncidentsByMonitorIDWithinRange(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time) ([]models.Incident, error) {
	const query = `
		SELECT i.id, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		WHERE im.monitor_id = $1
		  AND (i.started_at < $3 AND (i.resolved_at IS NULL OR i.resolved_at >= $2))
		ORDER BY i.started_at DESC, i.id DESC
	`

	var incidents []models.Incident
	if err := pgxscan.Select(ctx, tx, &incidents, query, monitorID, start, end); err != nil {
		return nil, err
	}

	return incidents, nil
}
