package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yorukot/kymarium/models"
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
			ping.RegionID,
			ping.Latency,
			ping.Status,
		})
	}

	// Use COPY for performance and lower lock contention during bursts.
	// COPY cannot do ON CONFLICT, so we use a savepoint to fall back to an upsert
	// when we hit duplicates (e.g., at-least-once retries or timestamp rounding).
	const savepoint = "batch_insert_pings_copy"
	if _, err := tx.Exec(ctx, "SAVEPOINT "+savepoint); err != nil {
		return err
	}

	copied, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pings"},
		[]string{"time", "monitor_id", "region_id", "latency", "status"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		if _, rbErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+savepoint); rbErr != nil {
			return fmt.Errorf("copyfrom pings failed: %w (rollback to savepoint failed: %v)", err, rbErr)
		}
		_, _ = tx.Exec(ctx, "RELEASE SAVEPOINT "+savepoint)

		if isUniqueViolation(err) {
			return upsertPings(ctx, tx, pings)
		}
		return err
	}

	_, _ = tx.Exec(ctx, "RELEASE SAVEPOINT "+savepoint)

	if copied != int64(len(rows)) {
		return fmt.Errorf("expected to copy %d rows, copied %d", len(rows), copied)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func upsertPings(ctx context.Context, tx pgx.Tx, pings []models.Ping) error {
	const query = `
		INSERT INTO pings (time, monitor_id, region_id, latency, status)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (time, monitor_id, region_id)
		DO UPDATE SET
			latency = EXCLUDED.latency,
			status = EXCLUDED.status
	`

	batch := &pgx.Batch{}
	for _, ping := range pings {
		batch.Queue(query, ping.Time, ping.MonitorID, ping.RegionID, ping.Latency, ping.Status)
	}

	results := tx.SendBatch(ctx, batch)

	for range pings {
		if _, err := results.Exec(); err != nil {
			_ = results.Close()
			return err
		}
	}

	return results.Close()
}

// ListRecentPingsByMonitorIDAndRegion fetches the latest pings for a monitor in a region, ordered newest first.
func (r *PGRepository) ListRecentPingsByMonitorIDAndRegion(ctx context.Context, tx pgx.Tx, monitorID int64, regionID int64, limit int) ([]models.Ping, error) {
	if limit <= 0 {
		return []models.Ping{}, nil
	}

	const query = `
		SELECT time, monitor_id, region_id, latency, status
		FROM pings
		WHERE monitor_id = $1 AND region_id = $2
		ORDER BY time DESC
		LIMIT $3
	`

	var pings []models.Ping
	if err := pgxscan.Select(ctx, tx, &pings, query, monitorID, regionID, limit); err != nil {
		return nil, err
	}

	return pings, nil
}
