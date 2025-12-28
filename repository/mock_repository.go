package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/yorukot/knocker/models"
)

// MockRepository is a testify-based mock implementing Repository for unit tests.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) StartTransaction(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	tx, _ := args.Get(0).(pgx.Tx)
	return tx, args.Error(1)
}

func (m *MockRepository) DeferRollback(tx pgx.Tx, ctx context.Context) {
	m.Called(tx, ctx)
}

func (m *MockRepository) CommitTransaction(tx pgx.Tx, ctx context.Context) error {
	args := m.Called(tx, ctx)
	return args.Error(0)
}

func (m *MockRepository) CreateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) error {
	args := m.Called(ctx, tx, statusPage)
	return args.Error(0)
}

func (m *MockRepository) UpdateStatusPage(ctx context.Context, tx pgx.Tx, statusPage models.StatusPage) (*models.StatusPage, error) {
	args := m.Called(ctx, tx, statusPage)
	page, _ := args.Get(0).(*models.StatusPage)
	return page, args.Error(1)
}

func (m *MockRepository) GetStatusPageByID(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) (*models.StatusPage, error) {
	args := m.Called(ctx, tx, teamID, statusPageID)
	page, _ := args.Get(0).(*models.StatusPage)
	return page, args.Error(1)
}

func (m *MockRepository) GetStatusPageBySlug(ctx context.Context, tx pgx.Tx, slug string) (*models.StatusPage, error) {
	args := m.Called(ctx, tx, slug)
	page, _ := args.Get(0).(*models.StatusPage)
	return page, args.Error(1)
}

func (m *MockRepository) ListStatusPagesByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.StatusPage, error) {
	args := m.Called(ctx, tx, teamID)
	pages, _ := args.Get(0).([]models.StatusPage)
	return pages, args.Error(1)
}

func (m *MockRepository) ListStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageGroup, error) {
	args := m.Called(ctx, tx, statusPageID)
	groups, _ := args.Get(0).([]models.StatusPageGroup)
	return groups, args.Error(1)
}

func (m *MockRepository) ListStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) ([]models.StatusPageMonitor, error) {
	args := m.Called(ctx, tx, statusPageID)
	monitors, _ := args.Get(0).([]models.StatusPageMonitor)
	return monitors, args.Error(1)
}

func (m *MockRepository) CreateStatusPageGroups(ctx context.Context, tx pgx.Tx, groups []models.StatusPageGroup) error {
	args := m.Called(ctx, tx, groups)
	return args.Error(0)
}

func (m *MockRepository) CreateStatusPageMonitors(ctx context.Context, tx pgx.Tx, monitors []models.StatusPageMonitor) error {
	args := m.Called(ctx, tx, monitors)
	return args.Error(0)
}

func (m *MockRepository) DeleteStatusPage(ctx context.Context, tx pgx.Tx, teamID, statusPageID int64) error {
	args := m.Called(ctx, tx, teamID, statusPageID)
	return args.Error(0)
}

func (m *MockRepository) DeleteStatusPageMonitorsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error {
	args := m.Called(ctx, tx, statusPageID)
	return args.Error(0)
}

func (m *MockRepository) DeleteStatusPageGroupsByStatusPageID(ctx context.Context, tx pgx.Tx, statusPageID int64) error {
	args := m.Called(ctx, tx, statusPageID)
	return args.Error(0)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.User, error) {
	args := m.Called(ctx, tx, email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*models.User, error) {
	args := m.Called(ctx, tx, userID)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockRepository) UpdateUserVerification(ctx context.Context, tx pgx.Tx, userID int64, verified bool, verifyCode *string, updatedAt time.Time) error {
	args := m.Called(ctx, tx, userID, verified, verifyCode, updatedAt)
	return args.Error(0)
}

func (m *MockRepository) GetAccountByEmail(ctx context.Context, tx pgx.Tx, email string) (*models.Account, error) {
	args := m.Called(ctx, tx, email)
	account, _ := args.Get(0).(*models.Account)
	return account, args.Error(1)
}

func (m *MockRepository) GetAccountWithUserByProviderUserID(ctx context.Context, tx pgx.Tx, provider models.Provider, providerUserID string) (*models.Account, *models.User, error) {
	args := m.Called(ctx, tx, provider, providerUserID)
	account, _ := args.Get(0).(*models.Account)
	user, _ := args.Get(1).(*models.User)
	return account, user, args.Error(2)
}

func (m *MockRepository) GetRefreshTokenByToken(ctx context.Context, tx pgx.Tx, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, tx, token)
	refreshToken, _ := args.Get(0).(*models.RefreshToken)
	return refreshToken, args.Error(1)
}

func (m *MockRepository) CreateAccount(ctx context.Context, tx pgx.Tx, account models.Account) error {
	args := m.Called(ctx, tx, account)
	return args.Error(0)
}

func (m *MockRepository) CreateUserAndAccount(ctx context.Context, tx pgx.Tx, user models.User, account models.Account) error {
	args := m.Called(ctx, tx, user, account)
	return args.Error(0)
}

func (m *MockRepository) CreateOAuthToken(ctx context.Context, tx pgx.Tx, oauthToken models.OAuthToken) error {
	args := m.Called(ctx, tx, oauthToken)
	return args.Error(0)
}

func (m *MockRepository) CreateRefreshToken(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	args := m.Called(ctx, tx, token)
	return args.Error(0)
}

func (m *MockRepository) UpdateRefreshTokenUsedAt(ctx context.Context, tx pgx.Tx, token models.RefreshToken) error {
	args := m.Called(ctx, tx, token)
	return args.Error(0)
}

func (m *MockRepository) ListTeamsByUserID(ctx context.Context, tx pgx.Tx, userID int64) ([]models.TeamWithRole, error) {
	args := m.Called(ctx, tx, userID)
	teams, _ := args.Get(0).([]models.TeamWithRole)
	return teams, args.Error(1)
}

func (m *MockRepository) GetTeamForUser(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamWithRole, error) {
	args := m.Called(ctx, tx, teamID, userID)
	team, _ := args.Get(0).(*models.TeamWithRole)
	return team, args.Error(1)
}

func (m *MockRepository) GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, teamID, userID int64) (*models.TeamMember, error) {
	args := m.Called(ctx, tx, teamID, userID)
	member, _ := args.Get(0).(*models.TeamMember)
	return member, args.Error(1)
}

func (m *MockRepository) CreateTeam(ctx context.Context, tx pgx.Tx, team models.Team) error {
	args := m.Called(ctx, tx, team)
	return args.Error(0)
}

func (m *MockRepository) CreateTeamMember(ctx context.Context, tx pgx.Tx, member models.TeamMember) error {
	args := m.Called(ctx, tx, member)
	return args.Error(0)
}

func (m *MockRepository) UpdateTeamName(ctx context.Context, tx pgx.Tx, teamID int64, name string, updatedAt time.Time) (*models.Team, error) {
	args := m.Called(ctx, tx, teamID, name, updatedAt)
	team, _ := args.Get(0).(*models.Team)
	return team, args.Error(1)
}

func (m *MockRepository) DeleteTeam(ctx context.Context, tx pgx.Tx, teamID int64) error {
	args := m.Called(ctx, tx, teamID)
	return args.Error(0)
}

func (m *MockRepository) ListNotificationsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Notification, error) {
	args := m.Called(ctx, tx, teamID)
	notifications, _ := args.Get(0).([]models.Notification)
	return notifications, args.Error(1)
}

func (m *MockRepository) GetNotificationByID(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) (*models.Notification, error) {
	args := m.Called(ctx, tx, teamID, notificationID)
	notification, _ := args.Get(0).(*models.Notification)
	return notification, args.Error(1)
}

func (m *MockRepository) CreateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) error {
	args := m.Called(ctx, tx, notification)
	return args.Error(0)
}

func (m *MockRepository) UpdateNotification(ctx context.Context, tx pgx.Tx, notification models.Notification) (*models.Notification, error) {
	args := m.Called(ctx, tx, notification)
	updated, _ := args.Get(0).(*models.Notification)
	return updated, args.Error(1)
}

func (m *MockRepository) DeleteNotification(ctx context.Context, tx pgx.Tx, teamID, notificationID int64) error {
	args := m.Called(ctx, tx, teamID, notificationID)
	return args.Error(0)
}

func (m *MockRepository) CreateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) error {
	args := m.Called(ctx, tx, monitor)
	return args.Error(0)
}

func (m *MockRepository) ListMonitorsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Monitor, error) {
	args := m.Called(ctx, tx, teamID)
	monitors, _ := args.Get(0).([]models.Monitor)
	return monitors, args.Error(1)
}

func (m *MockRepository) ListMonitorsByIDs(ctx context.Context, tx pgx.Tx, teamID int64, monitorIDs []int64) ([]models.Monitor, error) {
	args := m.Called(ctx, tx, teamID, monitorIDs)
	monitors, _ := args.Get(0).([]models.Monitor)
	return monitors, args.Error(1)
}

func (m *MockRepository) GetMonitorByID(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) (*models.Monitor, error) {
	args := m.Called(ctx, tx, teamID, monitorID)
	monitor, _ := args.Get(0).(*models.Monitor)
	return monitor, args.Error(1)
}

func (m *MockRepository) UpdateMonitor(ctx context.Context, tx pgx.Tx, monitor models.Monitor) (*models.Monitor, error) {
	args := m.Called(ctx, tx, monitor)
	updated, _ := args.Get(0).(*models.Monitor)
	return updated, args.Error(1)
}

func (m *MockRepository) DeleteMonitor(ctx context.Context, tx pgx.Tx, teamID, monitorID int64) error {
	args := m.Called(ctx, tx, teamID, monitorID)
	return args.Error(0)
}

func (m *MockRepository) ListMonitorsDueForCheck(ctx context.Context, tx pgx.Tx) ([]models.Monitor, error) {
	args := m.Called(ctx, tx)
	monitors, _ := args.Get(0).([]models.Monitor)
	return monitors, args.Error(1)
}

func (m *MockRepository) BatchUpdateMonitorsLastChecked(ctx context.Context, tx pgx.Tx, monitorIDs []int64, nextChecks []time.Time, lastChecked time.Time) error {
	args := m.Called(ctx, tx, monitorIDs, nextChecks, lastChecked)
	return args.Error(0)
}

func (m *MockRepository) ListRegionsByIDs(ctx context.Context, tx pgx.Tx, regionIDs []int64) ([]models.Region, error) {
	args := m.Called(ctx, tx, regionIDs)
	regions, _ := args.Get(0).([]models.Region)
	return regions, args.Error(1)
}

func (m *MockRepository) BatchInsertPings(ctx context.Context, tx pgx.Tx, pings []models.Ping) error {
	args := m.Called(ctx, tx, pings)
	return args.Error(0)
}

func (m *MockRepository) CreateMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64, notificationIDs []int64) error {
	args := m.Called(ctx, tx, monitorID, notificationIDs)
	return args.Error(0)
}

func (m *MockRepository) DeleteMonitorNotifications(ctx context.Context, tx pgx.Tx, monitorID int64) error {
	args := m.Called(ctx, tx, monitorID)
	return args.Error(0)
}

func (m *MockRepository) GetNotificationIDsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]int64, error) {
	args := m.Called(ctx, tx, monitorID)
	notificationIDs, _ := args.Get(0).([]int64)
	return notificationIDs, args.Error(1)
}

func (m *MockRepository) CreateMonitorRegions(ctx context.Context, tx pgx.Tx, monitorID int64, regions []models.Region) error {
	args := m.Called(ctx, tx, monitorID, regions)
	return args.Error(0)
}

func (m *MockRepository) DeleteMonitorRegions(ctx context.Context, tx pgx.Tx, monitorID int64) error {
	args := m.Called(ctx, tx, monitorID)
	return args.Error(0)
}

func (m *MockRepository) GetOpenIncidentByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) (*models.Incident, error) {
	args := m.Called(ctx, tx, monitorID)
	incident, _ := args.Get(0).(*models.Incident)
	return incident, args.Error(1)
}

func (m *MockRepository) CreateIncident(ctx context.Context, tx pgx.Tx, incident models.Incident) error {
	args := m.Called(ctx, tx, incident)
	return args.Error(0)
}

func (m *MockRepository) CreateIncidentMonitor(ctx context.Context, tx pgx.Tx, incidentID, monitorID int64) error {
	args := m.Called(ctx, tx, incidentID, monitorID)
	return args.Error(0)
}

func (m *MockRepository) MarkIncidentResolved(ctx context.Context, tx pgx.Tx, incidentID int64, resolvedAt, updatedAt time.Time) error {
	args := m.Called(ctx, tx, incidentID, resolvedAt, updatedAt)
	return args.Error(0)
}

func (m *MockRepository) CreateEventTimeline(ctx context.Context, tx pgx.Tx, timeline models.EventTimeline) error {
	args := m.Called(ctx, tx, timeline)
	return args.Error(0)
}

func (m *MockRepository) GetLastEventTimeline(ctx context.Context, tx pgx.Tx, incidentID int64) (*models.EventTimeline, error) {
	args := m.Called(ctx, tx, incidentID)
	event, _ := args.Get(0).(*models.EventTimeline)
	return event, args.Error(1)
}

func (m *MockRepository) ListIncidentsByMonitorID(ctx context.Context, tx pgx.Tx, monitorID int64) ([]models.Incident, error) {
	args := m.Called(ctx, tx, monitorID)
	incidents, _ := args.Get(0).([]models.Incident)
	return incidents, args.Error(1)
}

func (m *MockRepository) ListIncidentsByTeamID(ctx context.Context, tx pgx.Tx, teamID int64) ([]models.Incident, error) {
	args := m.Called(ctx, tx, teamID)
	incidents, _ := args.Get(0).([]models.Incident)
	return incidents, args.Error(1)
}

func (m *MockRepository) ListPublicIncidentsByMonitorIDs(ctx context.Context, tx pgx.Tx, monitorIDs []int64) ([]models.IncidentWithMonitorID, error) {
	args := m.Called(ctx, tx, monitorIDs)
	incidents, _ := args.Get(0).([]models.IncidentWithMonitorID)
	return incidents, args.Error(1)
}

func (m *MockRepository) ListPublicEventTimelinesByIncidentIDs(ctx context.Context, tx pgx.Tx, incidentIDs []int64) ([]models.EventTimeline, error) {
	args := m.Called(ctx, tx, incidentIDs)
	events, _ := args.Get(0).([]models.EventTimeline)
	return events, args.Error(1)
}

func (m *MockRepository) GetIncidentByID(ctx context.Context, tx pgx.Tx, monitorID, incidentID int64) (*models.Incident, error) {
	args := m.Called(ctx, tx, monitorID, incidentID)
	incident, _ := args.Get(0).(*models.Incident)
	return incident, args.Error(1)
}

func (m *MockRepository) GetIncidentByIDForTeam(ctx context.Context, tx pgx.Tx, teamID, incidentID int64) (*models.Incident, error) {
	args := m.Called(ctx, tx, teamID, incidentID)
	incident, _ := args.Get(0).(*models.Incident)
	return incident, args.Error(1)
}

func (m *MockRepository) ListEventTimelinesByIncidentID(ctx context.Context, tx pgx.Tx, incidentID int64) ([]models.EventTimeline, error) {
	args := m.Called(ctx, tx, incidentID)
	events, _ := args.Get(0).([]models.EventTimeline)
	return events, args.Error(1)
}

func (m *MockRepository) UpdateIncidentStatus(ctx context.Context, tx pgx.Tx, incidentID int64, status models.IncidentStatus, resolvedAt *time.Time, updatedAt time.Time) (*models.Incident, error) {
	args := m.Called(ctx, tx, incidentID, status, resolvedAt, updatedAt)
	incident, _ := args.Get(0).(*models.Incident)
	return incident, args.Error(1)
}

func (m *MockRepository) UpdateIncidentSettings(ctx context.Context, tx pgx.Tx, incidentID int64, isPublic bool, autoResolve bool, title *string, updatedAt time.Time) (*models.Incident, error) {
	args := m.Called(ctx, tx, incidentID, isPublic, autoResolve, title, updatedAt)
	incident, _ := args.Get(0).(*models.Incident)
	return incident, args.Error(1)
}

func (m *MockRepository) ListRecentPingsByMonitorIDAndRegion(ctx context.Context, tx pgx.Tx, monitorID int64, regionID int64, limit int) ([]models.Ping, error) {
	args := m.Called(ctx, tx, monitorID, regionID, limit)
	pings, _ := args.Get(0).([]models.Ping)
	return pings, args.Error(1)
}

func (m *MockRepository) UpdateMonitorStatus(ctx context.Context, tx pgx.Tx, monitorID int64, status models.MonitorStatus, updatedAt time.Time) error {
	args := m.Called(ctx, tx, monitorID, status, updatedAt)
	return args.Error(0)
}

func (m *MockRepository) ListAllRegions(ctx context.Context, tx pgx.Tx) ([]models.Region, error) {
	args := m.Called(ctx, tx)
	regions, _ := args.Get(0).([]models.Region)
	return regions, args.Error(1)
}

func (m *MockRepository) GetMonitorAnalytics(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time, regionID *int64) ([]models.MonitorAnalyticsBucket, error) {
	args := m.Called(ctx, tx, monitorID, start, end, regionID)
	buckets, _ := args.Get(0).([]models.MonitorAnalyticsBucket)
	return buckets, args.Error(1)
}

func (m *MockRepository) ListMonitorDailySummaryByMonitorIDs(ctx context.Context, tx pgx.Tx, monitorIDs []int64, start time.Time, end time.Time) ([]models.MonitorDailySummary, error) {
	args := m.Called(ctx, tx, monitorIDs, start, end)
	summaries, _ := args.Get(0).([]models.MonitorDailySummary)
	return summaries, args.Error(1)
}

func (m *MockRepository) ListIncidentsByMonitorIDWithinRange(ctx context.Context, tx pgx.Tx, monitorID int64, start time.Time, end time.Time) ([]models.Incident, error) {
	args := m.Called(ctx, tx, monitorID, start, end)
	incidents, _ := args.Get(0).([]models.Incident)
	return incidents, args.Error(1)
}
