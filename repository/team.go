package repository

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// ListTeamsByUserID returns all teams the user is a member of.
func ListTeamsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.TeamWithRole, error) {
	query := `
		SELECT t.id, t.name, t.updated_at, t.created_at, tm.role
		FROM teams t
		INNER JOIN team_members tm ON tm.team_id = t.id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`

	var teams []models.TeamWithRole
	if err := pgxscan.Select(ctx, tx, &teams, query, userID); err != nil {
		return nil, err
	}

	return teams, nil
}

// GetTeamForUser returns the team if the user is a member of it.
func GetTeamForUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamWithRole, error) {
	query := `
		SELECT t.id, t.name, t.updated_at, t.created_at, tm.role
		FROM teams t
		INNER JOIN team_members tm ON tm.team_id = t.id
		WHERE t.id = $1 AND tm.user_id = $2
		LIMIT 1
	`

	var team models.TeamWithRole
	if err := tx.QueryRow(ctx, query, teamID, userID).Scan(
		&team.ID,
		&team.Name,
		&team.UpdatedAt,
		&team.CreatedAt,
		&team.Role,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

// GetTeamMemberByUserID returns the membership record for the user within the specified team.
func GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamMember, error) {
	query := `
		SELECT id, team_id, user_id, role, updated_at, created_at
		FROM team_members
		WHERE team_id = $1 AND user_id = $2
		LIMIT 1
	`

	var member models.TeamMember
	if err := tx.QueryRow(ctx, query, teamID, userID).Scan(
		&member.ID,
		&member.TeamID,
		&member.UserID,
		&member.Role,
		&member.UpdatedAt,
		&member.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}

// CreateTeam inserts a new team record.
func CreateTeam(ctx context.Context, tx pgx.Tx, team models.Team) error {
	query := `
		INSERT INTO teams (id, name, updated_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, team.ID, team.Name, team.UpdatedAt, team.CreatedAt)
	return err
}

// CreateTeamMember inserts a new team membership record.
func CreateTeamMember(ctx context.Context, tx pgx.Tx, member models.TeamMember) error {
	query := `
		INSERT INTO team_members (id, team_id, user_id, role, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := tx.Exec(ctx, query, member.ID, member.TeamID, member.UserID, member.Role, member.UpdatedAt, member.CreatedAt)
	return err
}

// UpdateTeamName updates a team's name.
func UpdateTeamName(ctx context.Context, tx pgx.Tx, teamID int64, name string, updatedAt time.Time) (*models.Team, error) {
	query := `
		UPDATE teams
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, updated_at, created_at
	`

	var team models.Team
	if err := tx.QueryRow(ctx, query, name, updatedAt, teamID).Scan(
		&team.ID,
		&team.Name,
		&team.UpdatedAt,
		&team.CreatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &team, nil
}

// DeleteTeam removes a team and its direct membership/invite records.
func DeleteTeam(ctx context.Context, tx pgx.Tx, teamID int64) error {
	if _, err := tx.Exec(ctx, `DELETE FROM team_invites WHERE team_id = $1`, teamID); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM team_members WHERE team_id = $1`, teamID); err != nil {
		return err
	}

	cmd, err := tx.Exec(ctx, `DELETE FROM teams WHERE id = $1`, teamID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
