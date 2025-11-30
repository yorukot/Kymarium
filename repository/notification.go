package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// CreateNotification inserts a notification record.
func CreateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) error {
	query := `
		INSERT INTO notifications (id, team_id, type, name, config, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(ctx, query,
		notification.ID,
		notification.TeamID,
		notification.Type,
		notification.Name,
		notification.Config,
		notification.UpdatedAt,
		notification.CreatedAt,
	)
	return err
}

// ListNotificationsByTeamID returns notifications belonging to a team.
func ListNotificationsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Notification, error) {
	query := `
		SELECT id, team_id, type, name, config, updated_at, created_at
		FROM notifications
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	var notifications []models.Notification
	if err := pgxscan.Select(ctx, tx, &notifications, query, teamID); err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetNotificationByID fetches a notification ensuring it belongs to the provided team.
func GetNotificationByID(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) (*models.Notification, error) {
	query := `
		SELECT id, team_id, type, name, config, updated_at, created_at
		FROM notifications
		WHERE id = $1 AND team_id = $2
	`

	var notification models.Notification
	if err := tx.QueryRow(ctx, query, notificationID, teamID).Scan(
		&notification.ID,
		&notification.TeamID,
		&notification.Type,
		&notification.Name,
		&notification.Config,
		&notification.UpdatedAt,
		&notification.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &notification, nil
}

// UpdateNotification updates a notification and returns the persisted record.
func UpdateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) (*models.Notification, error) {
	query := `
		UPDATE notifications
		SET type = $1, name = $2, config = $3, updated_at = $4
		WHERE id = $5 AND team_id = $6
		RETURNING id, team_id, type, name, config, updated_at, created_at
	`

	var updated models.Notification
	if err := tx.QueryRow(ctx, query,
		notification.Type,
		notification.Name,
		notification.Config,
		notification.UpdatedAt,
		notification.ID,
		notification.TeamID,
	).Scan(
		&updated.ID,
		&updated.TeamID,
		&updated.Type,
		&updated.Name,
		&updated.Config,
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

// DeleteNotification removes a notification belonging to a team.
func DeleteNotification(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) error {
	result, err := tx.Exec(ctx, `DELETE FROM notifications WHERE id = $1 AND team_id = $2`, notificationID, teamID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
