package repository

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/kymarium/models"
)

// GetUserByEmail retrieves a user by email address (through the accounts table)
func (r *PGRepository) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	query := `
		SELECT u.id, u.password_hash, u.display_name, u.avatar, u.verified, u.verify_code, u.created_at, u.updated_at
		FROM users u
		JOIN accounts a ON u.id = a.user_id
		WHERE a.email = $1 AND a.provider = $2
		LIMIT 1`

	var user models.User
	err := tx.QueryRow(ctx, query, email, models.ProviderEmail).Scan(
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
		return nil, nil // Not an error, just not found
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAccountByEmail retrieves an account by email address
func (r *PGRepository) GetAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.Account, error) {
	query := `SELECT id, provider, provider_user_id, user_id, email, is_primary, created_at, updated_at
	          FROM accounts
	          WHERE email = $1
	          LIMIT 1`

	var account models.Account
	err := tx.QueryRow(ctx, query, email).Scan(
		&account.ID,
		&account.Provider,
		&account.ProviderUserID,
		&account.UserID,
		&account.Email,
		&account.IsPrimary,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // Not an error, just not found
	}

	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetUserIDByEmail retrieves the distinct user ID for an email if it exists.
func (r *PGRepository) GetUserIDByEmail(ctx context.Context, tx pgx.Tx, email string) (*int64, error) {
	query := `SELECT DISTINCT user_id FROM accounts WHERE email = $1`

	var userIDs []int64
	if err := pgxscan.Select(ctx, tx, &userIDs, query, email); err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return nil, nil
	}

	if len(userIDs) > 1 {
		return nil, errors.New("email belongs to multiple users")
	}

	return &userIDs[0], nil
}

// GetAccountWithUserByProviderUserID retrieves the account and its associated user
func (r *PGRepository) GetAccountWithUserByProviderUserID(ctx context.Context, db pgx.Tx, provider models.Provider, providerUserID string) (*models.Account, *models.User, error) {
	query := `
		SELECT
			a.id AS "a.id", a.provider AS "a.provider", a.provider_user_id AS "a.provider_user_id", a.user_id AS "a.user_id", a.is_primary AS "a.is_primary",
			u.id AS "u.id", u.verified AS "u.verified", u.verify_code AS "u.verify_code", u.created_at AS "u.created_at", u.updated_at AS "u.updated_at"
		FROM accounts a
		JOIN users u ON a.user_id = u.id
		WHERE a.provider = $1 AND a.provider_user_id = $2
	`

	// Using aliases to scan into both Account and User
	var result struct {
		A models.Account `db:"a"`
		U models.User    `db:"u"`
	}

	err := pgxscan.Get(ctx, db, &result, query, provider, providerUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	} else if err != nil {
		return nil, nil, err
	}

	return &result.A, &result.U, nil
}

// GetSessionByToken retrieves a session by its token value
func (r *PGRepository) GetSessionByToken(ctx context.Context, tx pgx.Tx, token string) (*models.Session, error) {
	query := `SELECT id, user_id, token, user_agent, ip, expires_at, created_at
	          FROM sessions
	          WHERE token = $1
	          LIMIT 1`

	var session models.Session
	err := tx.QueryRow(ctx, query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.UserAgent,
		&session.IP,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token value
func (r *PGRepository) GetRefreshTokenByToken(ctx context.Context, tx pgx.Tx, token string) (*models.RefreshToken, error) {
	query := `SELECT id, user_id, token, user_agent, ip, used_at, created_at
	          FROM refresh_tokens
	          WHERE token = $1
	          LIMIT 1`

	var refreshToken models.RefreshToken
	err := tx.QueryRow(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.UserAgent,
		&refreshToken.IP,
		&refreshToken.UsedAt,
		&refreshToken.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // Not an error, just not found
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// CreateAccount creates a new account
func (r *PGRepository) CreateAccount(ctx context.Context, tx pgx.Tx, account models.Account) error {
	query := `INSERT INTO accounts (id, provider, provider_user_id, user_id, email, is_primary, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(ctx, query,
		account.ID,
		account.Provider,
		account.ProviderUserID,
		account.UserID,
		account.Email,
		account.IsPrimary,
		account.CreatedAt,
		account.UpdatedAt,
	)

	return err
}

// CreateUserAndAccount creates a new user and associated account in a transaction
func (r *PGRepository) CreateUserAndAccount(ctx context.Context, tx pgx.Tx, user models.User, account models.Account) error {
	// Insert user
	userQuery := `INSERT INTO users (id, password_hash, display_name, avatar, verified, verify_code, created_at, updated_at)
	              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.Exec(ctx, userQuery,
		user.ID,
		user.PasswordHash,
		user.DisplayName,
		user.Avatar,
		user.Verified,
		user.VerifyCode,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Insert account
	accountQuery := `INSERT INTO accounts (id, provider, provider_user_id, user_id, email, is_primary, created_at, updated_at)
	                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(ctx, accountQuery,
		account.ID,
		account.Provider,
		account.ProviderUserID,
		account.UserID,
		account.Email,
		account.IsPrimary,
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// CreateOAuthToken creates a new OAuth token
func (r *PGRepository) CreateOAuthToken(ctx context.Context, db pgx.Tx, oauthToken models.OAuthToken) error {
	query := `
		INSERT INTO oauth_tokens (
			account_id,
			access_token,
			refresh_token,
			expiry,
			token_type,
			provider,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (account_id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expiry = EXCLUDED.expiry,
			token_type = EXCLUDED.token_type,
			updated_at = EXCLUDED.updated_at
	`

	_, err := db.Exec(ctx,
		query,
		oauthToken.AccountID,
		oauthToken.AccessToken,
		oauthToken.RefreshToken,
		oauthToken.Expiry,
		oauthToken.TokenType,
		oauthToken.Provider,
		oauthToken.CreatedAt,
		oauthToken.UpdatedAt,
	)
	return err
}

// CreateSession creates a new session in the database
func (r *PGRepository) CreateSession(ctx context.Context, tx pgx.Tx, session models.Session) error {
	query := `INSERT INTO sessions (id, user_id, token, user_agent, ip, expires_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.Token,
		session.UserAgent,
		session.IP,
		session.ExpiresAt,
		session.CreatedAt,
	)

	return err
}

// CreateRefreshToken creates a new refresh token in the database
func (r *PGRepository) CreateRefreshToken(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, user_agent, ip, used_at, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := tx.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.UserAgent,
		token.IP,
		token.UsedAt,
		token.CreatedAt,
	)

	return err
}

// ListActiveSessionsByUserID returns active sessions for a user.
func (r *PGRepository) ListActiveSessionsByUserID(ctx context.Context, tx pgx.Tx, userID int64, now time.Time) ([]models.Session, error) {
	query := `SELECT id, user_id, token, user_agent, ip, expires_at, created_at
	          FROM sessions
	          WHERE user_id = $1 AND expires_at > $2
	          ORDER BY created_at DESC`

	var sessions []models.Session
	if err := pgxscan.Select(ctx, tx, &sessions, query, userID, now); err != nil {
		return nil, err
	}

	return sessions, nil
}

// DeleteSessionByID deletes a session by ID for a user.
func (r *PGRepository) DeleteSessionByID(ctx context.Context, tx pgx.Tx, userID, sessionID int64) (bool, error) {
	query := `DELETE FROM sessions WHERE id = $1 AND user_id = $2`
	result, err := tx.Exec(ctx, query, sessionID, userID)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

// DeleteSessionByToken deletes a session by token.
func (r *PGRepository) DeleteSessionByToken(ctx context.Context, tx pgx.Tx, token string) (bool, error) {
	query := `DELETE FROM sessions WHERE token = $1`
	result, err := tx.Exec(ctx, query, token)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

// DeleteSessionsExceptToken deletes all sessions except the one with the provided token.
func (r *PGRepository) DeleteSessionsExceptToken(ctx context.Context, tx pgx.Tx, userID int64, token string) (int64, error) {
	query := `DELETE FROM sessions WHERE user_id = $1 AND token != $2`
	result, err := tx.Exec(ctx, query, userID, token)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// UpdateRefreshTokenUsedAt updates the used_at timestamp for a refresh token
func (r *PGRepository) UpdateRefreshTokenUsedAt(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	query := `UPDATE refresh_tokens
	          SET used_at = $1
	          WHERE id = $2`

	_, err := tx.Exec(ctx, query, token.UsedAt, token.ID)
	return err
}

// ListAccountsByUserID retrieves all accounts for a user.
func (r *PGRepository) ListAccountsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.Account, error) {
	query := `SELECT id, provider, provider_user_id, user_id, email, is_primary, created_at, updated_at
	          FROM accounts
	          WHERE user_id = $1
	          ORDER BY created_at DESC`

	var accounts []models.Account
	if err := pgxscan.Select(ctx, tx, &accounts, query, userID); err != nil {
		return nil, err
	}

	return accounts, nil
}

// ListActiveRefreshTokensByUserID returns active (unused) refresh tokens for a user.
func (r *PGRepository) ListActiveRefreshTokensByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.RefreshToken, error) {
	query := `SELECT id, user_id, token, user_agent, ip, used_at, created_at
	          FROM refresh_tokens
	          WHERE user_id = $1 AND used_at IS NULL
	          ORDER BY created_at DESC`

	var tokens []models.RefreshToken
	if err := pgxscan.Select(ctx, tx, &tokens, query, userID); err != nil {
		return nil, err
	}

	return tokens, nil
}

// UpdateRefreshTokenUsedAtByID marks a refresh token as used for a user.
func (r *PGRepository) UpdateRefreshTokenUsedAtByID(ctx context.Context, tx pgx.Tx, userID, tokenID int64, usedAt time.Time) (bool, error) {
	query := `UPDATE refresh_tokens
	          SET used_at = $1
	          WHERE id = $2 AND user_id = $3 AND used_at IS NULL`

	result, err := tx.Exec(ctx, query, usedAt, tokenID, userID)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

// UpdateRefreshTokensUsedAtExcept marks all refresh tokens as used except the provided token.
func (r *PGRepository) UpdateRefreshTokensUsedAtExcept(ctx context.Context, tx pgx.Tx, userID, tokenID int64, usedAt time.Time) (int64, error) {
	query := `UPDATE refresh_tokens
	          SET used_at = $1
	          WHERE user_id = $2 AND used_at IS NULL AND id != $3`

	result, err := tx.Exec(ctx, query, usedAt, userID, tokenID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
