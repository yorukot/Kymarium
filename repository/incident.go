package repository

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/id"
)

// GetOpenIncidentByMonitorID fetches the latest non-resolved incident for a monitor, if any.
func (r *PGRepository) GetOpenIncidentByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) (*models.Incident, error) {
	const query = `
		SELECT i.id, i.title, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		WHERE im.monitor_id = $1
		  AND i.status <> 'resolved'
		ORDER BY i.started_at DESC, i.id DESC
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

// ListPublicIncidentsByMonitorIDs returns public incidents for the provided monitors.
func (r *PGRepository) ListPublicIncidentsByMonitorIDs(ctx context.Context, tx pgx.Tx, monitorIDs []int64) ([]models.IncidentWithMonitorID, error) {
	if len(monitorIDs) == 0 {
		return []models.IncidentWithMonitorID{}, nil
	}

	const query = `
		SELECT
			i.id,
			i.title,
			i.status,
			i.severity,
			i.is_public,
			i.auto_resolve,
			i.started_at,
			i.resolved_at,
			i.created_at,
			i.updated_at,
			im.monitor_id
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		WHERE im.monitor_id = ANY($1)
		  AND i.is_public = true
		ORDER BY i.started_at DESC, i.id DESC
	`

	var incidents []models.IncidentWithMonitorID
	if err := pgxscan.Select(ctx, tx, &incidents, query, monitorIDs); err != nil {
		return nil, err
	}

	return incidents, nil
}

// ListPublicEventTimelinesByIncidentIDs returns human-authored timeline events for incidents.
func (r *PGRepository) ListPublicEventTimelinesByIncidentIDs(ctx context.Context, tx pgx.Tx, incidentIDs []int64) ([]models.EventTimeline, error) {
	if len(incidentIDs) == 0 {
		return []models.EventTimeline{}, nil
	}

	const query = `
		SELECT id, event_id, created_by, message, event_type, created_at, updated_at
		FROM event_timelines
		WHERE event_id = ANY($1)
		  AND created_by IS NOT NULL
		  AND message <> ''
		ORDER BY created_at ASC, id ASC
	`

	var events []models.EventTimeline
	if err := pgxscan.Select(ctx, tx, &events, query, incidentIDs); err != nil {
		return nil, err
	}

	return events, nil
}

// CreateIncident inserts a new incident row.
func (r *PGRepository) CreateIncident(ctx context.Context, tx pgx.Tx, incident models.Incident) error {
	const query = `
		INSERT INTO incidents (id, title, status, severity, is_public, auto_resolve, started_at, resolved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		incident.Title,
		incident.Status,
		incident.Severity,
		incident.IsPublic,
		incident.AutoResolve,
		incident.StartedAt,
		incident.ResolvedAt,
		incident.CreatedAt,
		incident.UpdatedAt,
	)
	return err
}

// CreateIncidentMonitor links an incident to a monitor.
func (r *PGRepository) CreateIncidentMonitor(ctx context.Context, tx pgx.Tx, incidentID, monitorID int64) error {
	const query = `
		INSERT INTO incident_monitors (id, incident_id, monitor_id)
		VALUES ($1, $2, $3)
	`

	junctionID, err := id.GetID()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, junctionID, incidentID, monitorID)
	return err
}

// MarkIncidentResolved closes an incident.
func (r *PGRepository) MarkIncidentResolved(ctx context.Context, tx pgx.Tx, incidentID int64, resolvedAt, updatedAt time.Time) error {
	const query = `
		UPDATE incidents
		SET status = 'resolved',
		    resolved_at = $2,
		    updated_at = $3
		WHERE id = $1
	`

	_, err := tx.Exec(ctx, query, incidentID, resolvedAt, updatedAt)
	return err
}

// CreateEventTimeline inserts an event timeline entry.
func (r *PGRepository) CreateEventTimeline(ctx context.Context, tx pgx.Tx, timeline models.EventTimeline) error {
	const query = `
		INSERT INTO event_timelines (id, event_id, created_by, message, event_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	idVal := timeline.ID
	if idVal == 0 {
		var err error
		idVal, err = id.GetID()
		if err != nil {
			return err
		}
	}

	timeline.ID = idVal

	_, err := tx.Exec(ctx, query,
		timeline.ID,
		timeline.IncidentID,
		timeline.CreatedBy,
		timeline.Message,
		timeline.EventType,
		timeline.CreatedAt,
		timeline.UpdatedAt,
	)
	return err
}

// GetLastEventTimeline returns the most recent timeline entry for an incident.
func (r *PGRepository) GetLastEventTimeline(ctx context.Context, tx pgx.Tx, incidentID int64) (*models.EventTimeline, error) {
	const query = `
		SELECT id, event_id, created_by, message, event_type, created_at, updated_at
		FROM event_timelines
		WHERE event_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`

	var event models.EventTimeline
	if err := pgxscan.Get(ctx, tx, &event, query, incidentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

// ListIncidentsByMonitorID returns all incidents for a monitor.
func (r *PGRepository) ListIncidentsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]models.Incident, error) {
	const query = `
		SELECT i.id, i.title, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		WHERE im.monitor_id = $1
		ORDER BY i.started_at DESC, i.id DESC
	`

	var incidents []models.Incident
	if err := pgxscan.Select(ctx, tx, &incidents, query, monitorID); err != nil {
		return nil, err
	}

	return incidents, nil
}

// GetIncidentByID fetches an incident scoped to the given monitor.
func (r *PGRepository) GetIncidentByID(ctx context.Context, tx pgx.Tx, monitorID, incidentID int64) (*models.Incident, error) {
	const query = `
		SELECT i.id, i.title, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		WHERE i.id = $1 AND im.monitor_id = $2
	`

	var incident models.Incident
	if err := pgxscan.Get(ctx, tx, &incident, query, incidentID, monitorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &incident, nil
}

// ListIncidentsByTeamID returns all incidents for a team via monitor membership.
func (r *PGRepository) ListIncidentsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Incident, error) {
	const query = `
		SELECT DISTINCT i.id, i.title, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		INNER JOIN monitors m ON m.id = im.monitor_id
		WHERE m.team_id = $1
		ORDER BY i.started_at DESC, i.id DESC
	`

	var incidents []models.Incident
	if err := pgxscan.Select(ctx, tx, &incidents, query, teamID); err != nil {
		return nil, err
	}
	return incidents, nil
}

// GetIncidentByIDForTeam fetches an incident ensuring it belongs to the team via monitor association.
func (r *PGRepository) GetIncidentByIDForTeam(ctx context.Context, tx pgx.Tx, teamID, incidentID int64) (*models.Incident, error) {
	const query = `
		SELECT i.id, i.title, i.status, i.severity, i.is_public, i.auto_resolve, i.started_at, i.resolved_at, i.created_at, i.updated_at
		FROM incidents i
		INNER JOIN incident_monitors im ON im.incident_id = i.id
		INNER JOIN monitors m ON m.id = im.monitor_id
		WHERE i.id = $1 AND m.team_id = $2
		LIMIT 1
	`

	var incident models.Incident
	if err := pgxscan.Get(ctx, tx, &incident, query, incidentID, teamID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &incident, nil
}

// ListEventTimelinesByIncidentID fetches all events for an incident in chronological order.
func (r *PGRepository) ListEventTimelinesByIncidentID(ctx context.Context, tx pgx.Tx, incidentID int64) ([]models.EventTimeline, error) {
	const query = `
		SELECT id, event_id, created_by, message, event_type, created_at, updated_at
		FROM event_timelines
		WHERE event_id = $1
		ORDER BY created_at ASC, id ASC
	`

	var events []models.EventTimeline
	if err := pgxscan.Select(ctx, tx, &events, query, incidentID); err != nil {
		return nil, err
	}

	return events, nil
}

// UpdateIncidentStatus updates the status (and optional resolved time) for an incident and returns the updated row.
func (r *PGRepository) UpdateIncidentStatus(ctx context.Context, tx pgx.Tx, incidentID int64, status models.IncidentStatus, resolvedAt *time.Time, updatedAt time.Time) (*models.Incident, error) {
	const query = `
		UPDATE incidents
		SET status = $2,
		    resolved_at = $3,
		    updated_at = $4
		WHERE id = $1
		RETURNING id, title, status, severity, is_public, auto_resolve, started_at, resolved_at, created_at, updated_at
	`

	var incident models.Incident
	if err := pgxscan.Get(ctx, tx, &incident, query, incidentID, status, resolvedAt, updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &incident, nil
}

// UpdateIncidentSettings updates visibility, auto-resolve flags, and title for an incident.
func (r *PGRepository) UpdateIncidentSettings(ctx context.Context, tx pgx.Tx, incidentID int64, isPublic bool, autoResolve bool, title *string, updatedAt time.Time) (*models.Incident, error) {
	const query = `
		UPDATE incidents
		SET is_public = $2,
		    auto_resolve = $3,
		    title = $4,
		    updated_at = $5
		WHERE id = $1
		RETURNING id, title, status, severity, is_public, auto_resolve, started_at, resolved_at, created_at, updated_at
	`

	var incident models.Incident
	if err := pgxscan.Get(ctx, tx, &incident, query, incidentID, isPublic, autoResolve, title, updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &incident, nil
}
