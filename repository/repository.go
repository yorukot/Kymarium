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

	// Status pages
	CreateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) error
	UpdateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) (*models.StatusPage, error)
	GetStatusPageByID(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) (*models.StatusPage, error)
	GetStatusPageBySlug(ctx context.Context, tx pgx.Tx, slug string) (*models.StatusPage, error)
	ListStatusPagesByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.StatusPage, error)
	ListStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageGroup, error)
	ListStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageMonitor, error)
	ListMonitorsByIDs(ctx context.Context, tx pgx.Tx, teamID int64, monitorIDs []int64) ([]models.Monitor, error)
	CreateStatusPageGroups(ctx context.Context, tx pgx.Tx, groups []models.StatusPageGroup) error
	CreateStatusPageMonitors(ctx context.Context, tx pgx.Tx, monitors []models.StatusPageMonitor) error
	DeleteStatusPage(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) error
	DeleteStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error
	DeleteStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error

	// Auth
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error)
	GetAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.Account, error)
	GetUserIDByEmail(ctx context.Context, tx pgx.Tx, email string) (*int64, error)
	GetAccountWithUserByProviderUserID(ctx context.Context, tx pgx.Tx, provider models.Provider, providerUserID string) (*models.Account, *models.User, error)
	GetRefreshTokenByToken(ctx context.Context, tx pgx.Tx, token string) (*models.RefreshToken, error)
	CreateAccount(ctx context.Context, tx pgx.Tx, account models.Account) error
	CreateUserAndAccount(ctx context.Context, tx pgx.Tx, user models.User, account models.Account) error
	CreateOAuthToken(ctx context.Context, tx pgx.Tx, oauthToken models.OAuthToken) error
	CreateRefreshToken(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error
	UpdateRefreshTokenUsedAt(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error

	// Users
	GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*models.User, error)
	UpdateUserVerification(ctx context.Context, tx pgx.Tx, userID int64, verified bool, verifyCode *string, updatedAt time.Time) error

	// Teams
	ListTeamsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.TeamWithRole, error)
	GetTeamForUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamWithRole, error)
	GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamMember, error)
	CreateTeam(ctx context.Context, tx pgx.Tx, team models.Team) error
	CreateTeamMember(ctx context.Context, tx pgx.Tx, member models.TeamMember) error
	DeleteTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) error
	UpdateTeamName(ctx context.Context, tx pgx.Tx, teamID int64, name string, updatedAt time.Time) (*models.Team, error)
	DeleteTeam(ctx context.Context, tx pgx.Tx, teamID int64) error

	// Team Invites
	CreateTeamInvite(ctx context.Context, tx pgx.Tx, invite models.TeamInvite) error
	GetTeamInviteByID(ctx context.Context, tx pgx.Tx, teamID, inviteID int64) (*models.TeamInvite, error)
	GetTeamInviteByToken(ctx context.Context, tx pgx.Tx, token string) (*models.TeamInvite, error)
	GetPendingTeamInviteByTeamAndUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamInvite, error)
	ListTeamInvitesByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.TeamInvite, error)
	UpdateTeamInviteStatus(ctx context.Context, tx pgx.Tx, inviteID int64, status models.InviteStatus, updatedAt time.Time, acceptedAt, rejectedAt, canceledAt *time.Time) (*models.TeamInvite, error)

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
	ListRegionsByIDs(ctx context.Context, tx pgx.Tx, regionIDs []int64) ([]models.Region, error)

	// Monitor-Notification junction table
	CreateMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64, notificationIDs []int64) error
	DeleteMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64) error
	GetNotificationIDsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]int64, error)

	// Monitor-Region junction table
	CreateMonitorRegions(ctx context.Context, tx pgx.Tx, monitorID int64, regions []models.Region) error
	DeleteMonitorRegions(ctx context.Context, tx pgx.Tx, monitorID int64) error

	// Pings
	BatchInsertPings(ctx context.Context, tx pgx.Tx, pings []models.Ping) error

	// Regions
	ListAllRegions(ctx context.Context, tx pgx.Tx) ([]models.Region, error)

	// Incidents
	GetOpenIncidentByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) (*models.Incident, error)
	CreateIncident(ctx context.Context, tx pgx.Tx, incident models.Incident) error
	CreateIncidentMonitor(ctx context.Context, tx pgx.Tx, incidentID, monitorID int64) error
	MarkIncidentResolved(ctx context.Context, tx pgx.Tx, incidentID int64, resolvedAt, updatedAt time.Time) error
	CreateEventTimeline(ctx context.Context, tx pgx.Tx, timeline models.EventTimeline) error
	GetLastEventTimeline(ctx context.Context, tx pgx.Tx, incidentID int64) (*models.EventTimeline, error)
	ListIncidentsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]models.Incident, error)
	ListIncidentsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Incident, error)
	ListPublicIncidentsByMonitorIDs(ctx context.Context, tx pgx.Tx, monitorIDs []int64) ([]models.IncidentWithMonitorID, error)
	ListPublicEventTimelinesByIncidentIDs(ctx context.Context, tx pgx.Tx, incidentIDs []int64) ([]models.EventTimeline, error)
	GetIncidentByID(ctx context.Context, tx pgx.Tx, monitorID, incidentID int64) (*models.Incident, error)
	GetIncidentByIDForTeam(ctx context.Context, tx pgx.Tx, teamID, incidentID int64) (*models.Incident, error)
	ListEventTimelinesByIncidentID(ctx context.Context, tx pgx.Tx, incidentID int64) ([]models.EventTimeline, error)
	UpdateIncidentStatus(ctx context.Context, tx pgx.Tx, incidentID int64, status models.IncidentStatus, resolvedAt *time.Time, updatedAt time.Time) (*models.Incident, error)
	UpdateIncidentSettings(ctx context.Context, tx pgx.Tx, incidentID int64, isPublic bool, autoResolve bool, title *string, updatedAt time.Time) (*models.Incident, error)
	ListRecentPingsByMonitorIDAndRegion(ctx context.Context, tx pgx.Tx, monitorID int64, regionID int64, limit int) ([]models.Ping, error)
	UpdateMonitorStatus(ctx context.Context, tx pgx.Tx, monitorID int64, status models.MonitorStatus, updatedAt time.Time) error

	// Analytics
	GetMonitorAnalytics(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time, regionID *int64) ([]models.MonitorAnalyticsBucket, error)
	ListMonitorDailySummaryByMonitorIDs(ctx context.Context, tx pgx.Tx, monitorIDs []int64, start time.Time, end time.Time) ([]models.MonitorDailySummary, error)
	ListIncidentsByMonitorIDWithinRange(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time) ([]models.Incident, error)
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
