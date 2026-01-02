package repository

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/kymarium/models"
)

// CreateStatusPage inserts a new status page.
func (r *PGRepository) CreateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) error {
	query := `
		INSERT INTO status_pages (id, team_id, title, slug, icon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(ctx, query,
		statusPage.ID,
		statusPage.TeamID,
		statusPage.Title,
		statusPage.Slug,
		statusPage.Icon,
		statusPage.CreatedAt,
		statusPage.UpdatedAt,
	)

	return err
}

// UpdateStatusPage updates slug/icon for a status page and returns the row.
func (r *PGRepository) UpdateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) (*models.StatusPage, error) {
	query := `
		UPDATE status_pages
		SET title = $1, slug = $2, icon = $3, updated_at = $4
		WHERE id = $5 AND team_id = $6
		RETURNING id, team_id, title, slug, icon, created_at, updated_at
	`

	var updated models.StatusPage
	if err := tx.QueryRow(ctx, query,
		statusPage.Title,
		statusPage.Slug,
		statusPage.Icon,
		statusPage.UpdatedAt,
		statusPage.ID,
		statusPage.TeamID,
	).Scan(
		&updated.ID,
		&updated.TeamID,
		&updated.Title,
		&updated.Slug,
		&updated.Icon,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}

// GetStatusPageByID fetches a status page ensuring it belongs to the team.
func (r *PGRepository) GetStatusPageByID(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) (*models.StatusPage, error) {
	query := `
		SELECT id, team_id, title, slug, icon, created_at, updated_at
		FROM status_pages
		WHERE id = $1 AND team_id = $2
	`

	var statusPage models.StatusPage
	if err := pgxscan.Get(ctx, tx, &statusPage, query, statusPageID, teamID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &statusPage, nil
}

// GetStatusPageBySlug returns a status page matching the slug.
func (r *PGRepository) GetStatusPageBySlug(ctx context.Context, tx pgx.Tx, slug string) (*models.StatusPage, error) {
	query := `
		SELECT id, team_id, title, slug, icon, created_at, updated_at
		FROM status_pages
		WHERE slug = $1
	`

	var statusPage models.StatusPage
	if err := pgxscan.Get(ctx, tx, &statusPage, query, slug); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &statusPage, nil
}

// ListStatusPagesByTeamID returns status pages belonging to a team.
func (r *PGRepository) ListStatusPagesByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.StatusPage, error) {
	query := `
		SELECT id, team_id, title, slug, icon, created_at, updated_at
		FROM status_pages
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	var statusPages []models.StatusPage
	if err := pgxscan.Select(ctx, tx, &statusPages, query, teamID); err != nil {
		return nil, err
	}

	return statusPages, nil
}

// ListStatusPageGroupsByStatusPageID returns groups for a status page.
func (r *PGRepository) ListStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageGroup, error) {
	query := `
		SELECT id, status_page_id, name, type, sort_order
		FROM status_page_groups
		WHERE status_page_id = $1
		ORDER BY sort_order ASC
	`

	var groups []models.StatusPageGroup
	if err := pgxscan.Select(ctx, tx, &groups, query, statusPageID); err != nil {
		return nil, err
	}

	return groups, nil
}

// ListStatusPageMonitorsByStatusPageID returns monitors for a status page.
func (r *PGRepository) ListStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageMonitor, error) {
	query := `
		SELECT id, status_page_id, monitor_id, group_id, name, type, sort_order
		FROM status_page_monitors
		WHERE status_page_id = $1
		ORDER BY group_id NULLS FIRST, sort_order ASC
	`

	var monitors []models.StatusPageMonitor
	if err := pgxscan.Select(ctx, tx, &monitors, query, statusPageID); err != nil {
		return nil, err
	}

	return monitors, nil
}

// CreateStatusPageGroups bulk inserts groups for a status page.
func (r *PGRepository) CreateStatusPageGroups(ctx context.Context, tx pgx.Tx, groups []models.StatusPageGroup) error {
	if len(groups) == 0 {
		return nil
	}

	query := `
		INSERT INTO status_page_groups (id, status_page_id, name, type, sort_order)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, group := range groups {
		if _, err := tx.Exec(ctx, query,
			group.ID,
			group.StatusPageID,
			group.Name,
			group.Type,
			group.SortOrder,
		); err != nil {
			return err
		}
	}

	return nil
}

// CreateStatusPageMonitors bulk inserts monitors for a status page.
func (r *PGRepository) CreateStatusPageMonitors(ctx context.Context, tx pgx.Tx, monitors []models.StatusPageMonitor) error {
	if len(monitors) == 0 {
		return nil
	}

	query := `
		INSERT INTO status_page_monitors (id, status_page_id, monitor_id, group_id, name, type, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, monitor := range monitors {
		if _, err := tx.Exec(ctx, query,
			monitor.ID,
			monitor.StatusPageID,
			monitor.MonitorID,
			monitor.GroupID,
			monitor.Name,
			monitor.Type,
			monitor.SortOrder,
		); err != nil {
			return err
		}
	}

	return nil
}

// DeleteStatusPage removes a status page belonging to a team.
func (r *PGRepository) DeleteStatusPage(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) error {
	result, err := tx.Exec(ctx, `DELETE FROM status_pages WHERE id = $1 AND team_id = $2`, statusPageID, teamID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// DeleteStatusPageMonitorsByStatusPageID removes monitors tied to the status page.
func (r *PGRepository) DeleteStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error {
	_, err := tx.Exec(ctx, `DELETE FROM status_page_monitors WHERE status_page_id = $1`, statusPageID)
	return err
}

// DeleteStatusPageGroupsByStatusPageID removes groups tied to the status page.
func (r *PGRepository) DeleteStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error {
	_, err := tx.Exec(ctx, `DELETE FROM status_page_groups WHERE status_page_id = $1`, statusPageID)
	return err
}
