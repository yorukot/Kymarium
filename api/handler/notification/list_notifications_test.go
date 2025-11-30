package notification

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/knocker/internal/testutil"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

func TestListNotifications_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)
	mockRepo.On("ListNotificationsByTeamID", mock.Anything, mock.Anything, int64(1)).
		Return([]models.Notification{
			{ID: 5, TeamID: 1, Name: "Discord Webhook", Type: "discord"},
			{ID: 6, TeamID: 1, Name: "Telegram Bot", Type: "telegram"},
		}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.ListNotifications(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Notifications retrieved successfully", resp["message"])

	data, ok := resp["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 2)
}

func TestListNotifications_EmptyList(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("ListNotificationsByTeamID", mock.Anything, mock.Anything, int64(1)).
		Return([]models.Notification{}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.ListNotifications(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	data, ok := resp["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 0)
}

func TestListNotifications_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications", nil)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.ListNotifications(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestListNotifications_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.ListNotifications(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestListNotifications_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/invalid/notifications", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("invalid")

	err := h.ListNotifications(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
