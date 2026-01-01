package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/yorukot/knocker/models"
)

// GetUserByID retrieves a user by their ID.
func (r *PGRepository) GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*models.User, error) {
	query := `
		SELECT id, password_hash, display_name, avatar, verified, verify_code, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1`

	var user models.User
	err := tx.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Avatar,
		&user.Verified,
		&user.VerifyCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserVerification updates the user's verification status and code.
func (r *PGRepository) UpdateUserVerification(ctx context.Context, tx pgx.Tx, userID int64, verified bool, verifyCode *string, updatedAt time.Time) error {
	query := `
		UPDATE users
		SET verified = $2,
			verify_code = $3,
			updated_at = $4
		WHERE id = $1`

	_, err := tx.Exec(ctx, query, userID, verified, verifyCode, updatedAt)
	return err
}

// UpdateUserProfile updates display name and avatar for a user.
func (r *PGRepository) UpdateUserProfile(ctx context.Context, tx pgx.Tx, userID int64, displayName string, avatar *string, updatedAt time.Time) (*models.User, error) {
	query := `
		UPDATE users
		SET display_name = $2,
			avatar = $3,
			updated_at = $4
		WHERE id = $1
		RETURNING id, password_hash, display_name, avatar, verified, verify_code, created_at, updated_at`

	var user models.User
	err := tx.QueryRow(ctx, query, userID, displayName, avatar, updatedAt).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Avatar,
		&user.Verified,
		&user.VerifyCode,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserPasswordHash updates the password hash for a user.
func (r *PGRepository) UpdateUserPasswordHash(ctx context.Context, tx pgx.Tx, userID int64, passwordHash string, updatedAt time.Time) error {
	query := `
		UPDATE users
		SET password_hash = $2,
			updated_at = $3
		WHERE id = $1`

	_, err := tx.Exec(ctx, query, userID, passwordHash, updatedAt)
	return err
}
