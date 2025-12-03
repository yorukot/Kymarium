package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/id"
)

// CreateMonitor inserts a monitor record.
func (r *PGRepository) CreateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) error {
	query := `
		INSERT INTO monitors (id, team_id, name, type, interval, config, last_checked, next_check, failure_threshold, recovery_threshold, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := tx.Exec(ctx, query,
		monitor.ID,
		monitor.TeamID,
		monitor.Name,
		monitor.Type,
		monitor.Interval,
		monitor.Config,
		monitor.LastChecked,
		monitor.NextCheck,
		monitor.FailureThreshold,
		monitor.RecoveryThreshold,
		monitor.UpdatedAt,
		monitor.CreatedAt,
	)
	return err
}

// ListMonitorsByTeamID returns monitors belonging to a team.
func (r *PGRepository) ListMonitorsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Monitor, error) {
	query := `
		SELECT
			m.id,
			m.team_id,
			m.name,
			m.type,
			m.interval,
			m.config,
			m.last_checked,
			m.next_check,
			m.failure_threshold,
			m.recovery_threshold,
			m.updated_at,
			m.created_at,
			COALESCE(
				array_agg(mn.notification_id ORDER BY mn.id)
					FILTER (WHERE mn.notification_id IS NOT NULL),
				'{}'
			) AS notification_ids
		FROM monitors m
		LEFT JOIN monitor_notifications mn ON m.id = mn.monitor_id
		WHERE m.team_id = $1
		GROUP BY m.id, m.team_id, m.name, m.type, m.interval, m.config, m.last_checked, m.next_check, m.failure_threshold, m.recovery_threshold, m.updated_at, m.created_at
		ORDER BY m.created_at DESC
	`

	var monitors []models.Monitor
	if err := pgxscan.Select(ctx, tx, &monitors, query, teamID); err != nil {
		return nil, err
	}

	return monitors, nil
}

// GetMonitorByID fetches a monitor ensuring it belongs to the provided team.
func (r *PGRepository) GetMonitorByID(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) (*models.Monitor, error) {
	query := `
		SELECT
			m.id,
			m.team_id,
			m.name,
			m.type,
			m.interval,
			m.config,
			m.last_checked,
			m.next_check,
			m.failure_threshold,
			m.recovery_threshold,
			m.updated_at,
			m.created_at,
			COALESCE(
				array_agg(mn.notification_id ORDER BY mn.id)
					FILTER (WHERE mn.notification_id IS NOT NULL),
				'{}'
			) AS notification_ids
		FROM monitors m
		LEFT JOIN monitor_notifications mn ON m.id = mn.monitor_id
		WHERE m.id = $1 AND m.team_id = $2
		GROUP BY m.id, m.team_id, m.name, m.type, m.interval, m.config, m.last_checked, m.next_check, m.failure_threshold, m.recovery_threshold, m.updated_at, m.created_at
	`

	var monitor models.Monitor
	if err := pgxscan.Get(ctx, tx, &monitor, query, monitorID, teamID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &monitor, nil
}

// UpdateMonitor updates a monitor and returns the persisted record.
func (r *PGRepository) UpdateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) (*models.Monitor, error) {
	query := `
		UPDATE monitors
		SET name = $1, type = $2, interval = $3, config = $4, last_checked = $5, next_check = $6, failure_threshold = $7, recovery_threshold = $8, updated_at = $9
		WHERE id = $10 AND team_id = $11
		RETURNING id, team_id, name, type, interval, config, last_checked, next_check, failure_threshold, recovery_threshold, updated_at, created_at
	`

	var updated models.Monitor
	if err := tx.QueryRow(ctx, query,
		monitor.Name,
		monitor.Type,
		monitor.Interval,
		monitor.Config,
		monitor.LastChecked,
		monitor.NextCheck,
		monitor.FailureThreshold,
		monitor.RecoveryThreshold,
		monitor.UpdatedAt,
		monitor.ID,
		monitor.TeamID,
	).Scan(
		&updated.ID,
		&updated.TeamID,
		&updated.Name,
		&updated.Type,
		&updated.Interval,
		&updated.Config,
		&updated.LastChecked,
		&updated.NextCheck,
		&updated.FailureThreshold,
		&updated.RecoveryThreshold,
		&updated.UpdatedAt,
		&updated.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteMonitor removes a monitor belonging to a team.
func (r *PGRepository) DeleteMonitor(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) error {
	result, err := tx.Exec(ctx, `DELETE FROM monitors WHERE id = $1 AND team_id = $2`, monitorID, teamID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// ListMonitorsDueForCheck fetches all monitors where next_check <= now
func (r *PGRepository) ListMonitorsDueForCheck(ctx context.Context, tx pgx.Tx) ([]models.Monitor, error) {
	query := `
		SELECT id, team_id, name, type, interval, config, last_checked, next_check, failure_threshold, recovery_threshold, updated_at, created_at
		FROM monitors
		WHERE next_check <= NOW()
		ORDER BY next_check ASC
	`

	var monitors []models.Monitor
	if err := pgxscan.Select(ctx, tx, &monitors, query); err != nil {
		return nil, err
	}

	return monitors, nil
}

// BatchUpdateMonitorsLastChecked updates last_checked and next_check for multiple monitors
// Each monitor can have its own next_check time based on its interval
// Uses PostgreSQL's unnest to efficiently update multiple rows with different values
func (r *PGRepository) BatchUpdateMonitorsLastChecked(ctx context.Context, tx pgx.Tx, monitorIDs []int64, nextChecks []time.Time, lastChecked time.Time) error {
	if len(monitorIDs) == 0 {
		return nil
	}

	if len(monitorIDs) != len(nextChecks) {
		return fmt.Errorf("monitorIDs and nextChecks must have the same length")
	}

	query := `
		UPDATE monitors AS m
		SET
			last_checked = $1,
			next_check   = data.next_check
		FROM (
			SELECT
				unnest($2::bigint[])      AS id,
				unnest($3::timestamptz[]) AS next_check
		) AS data
		WHERE m.id = data.id
	`

	_, err := tx.Exec(ctx, query, lastChecked, monitorIDs, nextChecks)
	return err
}

// CreateMonitorNotifications creates associations between a monitor and notification IDs in the junction table.
func (r *PGRepository) CreateMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64, notificationIDs []int64) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	query := `
		INSERT INTO monitor_notifications (id, monitor_id, notification_id)
		VALUES ($1, $2, $3)
	`

	for _, notificationID := range notificationIDs {
		junctionID, err := id.GetID()
		if err != nil {
			return fmt.Errorf("failed to generate junction table ID: %w", err)
		}

		if _, err := tx.Exec(ctx, query, junctionID, monitorID, notificationID); err != nil {
			return err
		}
	}

	return nil
}

// DeleteMonitorNotifications removes all notification associations for a monitor.
func (r *PGRepository) DeleteMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64) error {
	_, err := tx.Exec(ctx, `DELETE FROM monitor_notifications WHERE monitor_id = $1`, monitorID)
	return err
}

// GetNotificationIDsByMonitorID fetches all notification IDs associated with a monitor.
func (r *PGRepository) GetNotificationIDsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]int64, error) {
	query := `
		SELECT notification_id
		FROM monitor_notifications
		WHERE monitor_id = $1
		ORDER BY id
	`

	var notificationIDs []int64
	rows, err := tx.Query(ctx, query, monitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notificationID int64
		if err := rows.Scan(&notificationID); err != nil {
			return nil, err
		}
		notificationIDs = append(notificationIDs, notificationID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notificationIDs, nil
}
