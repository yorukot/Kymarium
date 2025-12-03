package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/knocker/models"
)

// Repository defines the database operations used across the app.
// It enables mocking the data layer in tests.
type Repository interface {
	StartTransaction(ctx context.Context) (pgx.Tx, error)
	DeferRollback(tx pgx.Tx, ctx context.Context)
	CommitTransaction(tx pgx.Tx, ctx context.Context) error

	// Auth
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error)
	GetAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.Account, error)
	GetAccountWithUserByProviderUserID(ctx context.Context, tx pgx.Tx, provider models.Provider, providerUserID string) (*models.Account, *models.User, error)
	GetRefreshTokenByToken(ctx context.Context, tx pgx.Tx, token string) (*models.RefreshToken, error)
	CreateAccount(ctx context.Context, tx pgx.Tx, account models.Account) error
	CreateUserAndAccount(ctx context.Context, tx pgx.Tx, user models.User, account models.Account) error
	CreateOAuthToken(ctx context.Context, tx pgx.Tx, oauthToken models.OAuthToken) error
	CreateRefreshToken(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error
	UpdateRefreshTokenUsedAt(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error

	// Teams
	ListTeamsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.TeamWithRole, error)
	GetTeamForUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamWithRole, error)
	GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamMember, error)
	CreateTeam(ctx context.Context, tx pgx.Tx, team models.Team) error
	CreateTeamMember(ctx context.Context, tx pgx.Tx, member models.TeamMember) error
	UpdateTeamName(ctx context.Context, tx pgx.Tx, teamID int64, name string, updatedAt time.Time) (*models.Team, error)
	DeleteTeam(ctx context.Context, tx pgx.Tx, teamID int64) error

	// Notifications
	ListNotificationsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Notification, error)
	GetNotificationByID(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) (*models.Notification, error)
	CreateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) error
	UpdateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) (*models.Notification, error)
	DeleteNotification(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) error

	// Monitors
	CreateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) error
	ListMonitorsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Monitor, error)
	GetMonitorByID(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) (*models.Monitor, error)
	UpdateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) (*models.Monitor, error)
	DeleteMonitor(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) error
	ListMonitorsDueForCheck(ctx context.Context, tx pgx.Tx) ([]models.Monitor, error)
	BatchUpdateMonitorsLastChecked(ctx context.Context, tx pgx.Tx, monitorIDs []int64, nextChecks []time.Time, lastChecked time.Time) error

	// Monitor-Notification junction table
	CreateMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64, notificationIDs []int64) error
	DeleteMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64) error
	GetNotificationIDsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]int64, error)

	// Pings
	BatchInsertPings(ctx context.Context, tx pgx.Tx, pings []models.Ping) error
}

// PGRepository is the production repository backed by pgx.
type PGRepository struct {
	DB *pgxpool.Pool
}

// New creates a repository backed by pgx pool.
func New(db *pgxpool.Pool) Repository {
	return &PGRepository{
		DB: db,
	}
}
