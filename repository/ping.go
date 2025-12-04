package repository

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// BatchInsertPings efficiently inserts ping records using COPY FROM.
// The caller is responsible for managing the transaction lifecycle.
func (r *PGRepository) BatchInsertPings(ctx context.Context, tx pgx.Tx, pings []models.Ping) error {
	if len(pings) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(pings))
	for _, ping := range pings {
		rows = append(rows, []any{
			ping.Time,
			ping.MonitorID,
			ping.Region,
			ping.Latency,
			ping.Status,
		})
	}

	// Use COPY for performance and lower lock contention during bursts.
	copied, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pings"},
		[]string{"time", "monitor_id", "region", "latency", "status"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	if copied != int64(len(rows)) {
		return fmt.Errorf("expected to copy %d rows, copied %d", len(rows), copied)
	}

	return nil
}

// ListRecentPingsByMonitorIDAndRegion fetches the latest pings for a monitor in a region, ordered newest first.
func (r *PGRepository) ListRecentPingsByMonitorIDAndRegion(ctx context.Context, tx pgx.Tx, monitorID int64, region string, limit int) ([]models.Ping, error) {
	if limit <= 0 {
		return []models.Ping{}, nil
	}

	const query = `
		SELECT time, monitor_id, region, latency, status
		FROM pings
		WHERE monitor_id = $1 AND region = $2
		ORDER BY time DESC
		LIMIT $3
	`

	var pings []models.Ping
	if err := pgxscan.Select(ctx, tx, &pings, query, monitorID, region, limit); err != nil {
		return nil, err
	}

	return pings, nil
}
