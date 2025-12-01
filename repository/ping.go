package repository

import (
	"context"
	"fmt"

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
			ping.Data,
		})
	}

	// Use COPY for performance and lower lock contention during bursts.
	copied, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"pings"},
		[]string{"time", "monitor_id", "region", "latency", "status", "data"},
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
