package repository

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// CreateTeamInvite inserts a new team invite record.
func (r *PGRepository) CreateTeamInvite(ctx context.Context, tx pgx.Tx, invite models.TeamInvite) error {
	query := `
		INSERT INTO team_invites (
			id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := tx.Exec(ctx, query,
		invite.ID,
		invite.TeamID,
		invite.InvitedBy,
		invite.InvitedTo,
		invite.InvitedEmail,
		invite.Role,
		invite.Status,
		invite.Token,
		invite.ExpiresAt,
		invite.AcceptedAt,
		invite.RejectedAt,
		invite.CanceledAt,
		invite.UpdatedAt,
		invite.CreatedAt,
	)

	return err
}

// GetTeamInviteByID fetches a team invite by id within a team.
func (r *PGRepository) GetTeamInviteByID(ctx context.Context, tx pgx.Tx, teamID, inviteID int64) (*models.TeamInvite, error) {
	query := `
		SELECT id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
		FROM team_invites
		WHERE id = $1 AND team_id = $2
		LIMIT 1
	`

	var invite models.TeamInvite
	err := pgxscan.Get(ctx, tx, &invite, query, inviteID, teamID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

// GetTeamInviteByToken fetches a team invite by token.
func (r *PGRepository) GetTeamInviteByToken(ctx context.Context, tx pgx.Tx, token string) (*models.TeamInvite, error) {
	query := `
		SELECT id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
		FROM team_invites
		WHERE token = $1
		LIMIT 1
	`

	var invite models.TeamInvite
	err := pgxscan.Get(ctx, tx, &invite, query, token)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

// GetPendingTeamInviteByTeamAndUser fetches a pending invite for a user on a team.
func (r *PGRepository) GetPendingTeamInviteByTeamAndUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamInvite, error) {
	query := `
		SELECT id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
		FROM team_invites
		WHERE team_id = $1 AND invited_to = $2 AND status = $3
		ORDER BY created_at DESC
		LIMIT 1
	`

	var invite models.TeamInvite
	err := pgxscan.Get(ctx, tx, &invite, query, teamID, userID, models.InviteStatusPending)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

// ListTeamInvitesByTeamID returns all invites for a team.
func (r *PGRepository) ListTeamInvitesByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.TeamInvite, error) {
	query := `
		SELECT id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
		FROM team_invites
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	var invites []models.TeamInvite
	if err := pgxscan.Select(ctx, tx, &invites, query, teamID); err != nil {
		return nil, err
	}

	return invites, nil
}

// UpdateTeamInviteStatus updates invite status and timestamps.
func (r *PGRepository) UpdateTeamInviteStatus(ctx context.Context, tx pgx.Tx, inviteID int64, status models.InviteStatus, updatedAt time.Time, acceptedAt, rejectedAt, canceledAt *time.Time) (*models.TeamInvite, error) {
	query := `
		UPDATE team_invites
		SET status = $1, updated_at = $2, accepted_at = $3, rejected_at = $4, canceled_at = $5
		WHERE id = $6
		RETURNING id, team_id, invited_by, invited_to, invited_email, role, status, token, expires_at,
			accepted_at, rejected_at, canceled_at, updated_at, created_at
	`

	var invite models.TeamInvite
	if err := tx.QueryRow(ctx, query, status, updatedAt, acceptedAt, rejectedAt, canceledAt, inviteID).Scan(
		&invite.ID,
		&invite.TeamID,
		&invite.InvitedBy,
		&invite.InvitedTo,
		&invite.InvitedEmail,
		&invite.Role,
		&invite.Status,
		&invite.Token,
		&invite.ExpiresAt,
		&invite.AcceptedAt,
		&invite.RejectedAt,
		&invite.CanceledAt,
		&invite.UpdatedAt,
		&invite.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &invite, nil
}
