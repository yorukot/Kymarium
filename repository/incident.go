package repository

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/utils/id"
)

// GetOpenIncidentByMonitorID fetches the latest non-resolved incident for a monitor, if any.
func (r *PGRepository) GetOpenIncidentByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) (*models.Incident, error) {
	const query = `
		SELECT id, monitor_id, status, started_at, resloved_at, created_at, updated_at
		FROM incidents
		WHERE monitor_id = $1
		  AND status <> 'resolved'
		ORDER BY started_at DESC
		LIMIT 1
	`

	var incident models.Incident
	if err := pgxscan.Get(ctx, tx, &incident, query, monitorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &incident, nil
}

// CreateIncident inserts a new incident row.
func (r *PGRepository) CreateIncident(ctx context.Context, tx pgx.Tx, incident models.Incident) error {
	const query = `
		INSERT INTO incidents (id, monitor_id, status, started_at, resloved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	idVal := incident.ID
	if idVal == 0 {
		var err error
		idVal, err = id.GetID()
		if err != nil {
			return err
		}
	}

	incident.ID = idVal

	_, err := tx.Exec(ctx, query,
		incident.ID,
		incident.MonitorID,
		incident.Status,
		incident.StartedAt,
		incident.ResolvedAt,
		incident.CreatedAt,
		incident.UpdatedAt,
	)
	return err
}

// MarkIncidentResolved closes an incident.
func (r *PGRepository) MarkIncidentResolved(ctx context.Context, tx pgx.Tx, incidentID int64, resolvedAt, updatedAt time.Time) error {
	const query = `
		UPDATE incidents
		SET status = 'resolved',
		    resloved_at = $2,
		    updated_at = $3
		WHERE id = $1
	`

	_, err := tx.Exec(ctx, query, incidentID, resolvedAt, updatedAt)
	return err
}

// CreateIncidentEvent inserts an incident event.
func (r *PGRepository) CreateIncidentEvent(ctx context.Context, tx pgx.Tx, event models.IncidentEvent) error {
	const query = `
		INSERT INTO incident_events (id, incident_id, created_by, message, event_type, public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	idVal := event.ID
	if idVal == 0 {
		var err error
		idVal, err = id.GetID()
		if err != nil {
			return err
		}
	}

	event.ID = idVal

	_, err := tx.Exec(ctx, query,
		event.ID,
		event.IncidentID,
		nil, // created_by is optional; set by API when available.
		event.Message,
		event.EventType,
		event.Public,
		event.CreatedAt,
		event.UpdatedAt,
	)
	return err
}

// GetLastIncidentEvent returns the most recent event for an incident.
func (r *PGRepository) GetLastIncidentEvent(ctx context.Context, tx pgx.Tx, incidentID int64) (*models.IncidentEvent, error) {
	const query = `
		SELECT id, incident_id, message, event_type, public, created_at, updated_at
		FROM incident_events
		WHERE incident_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`

	var event models.IncidentEvent
	if err := pgxscan.Get(ctx, tx, &event, query, incidentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}
