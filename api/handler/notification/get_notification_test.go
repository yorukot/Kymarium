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

func TestGetNotification_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "Test Notification", Type: "discord"}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.GetNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Notification retrieved successfully", resp["message"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Test Notification", data["name"])
}

func TestGetNotification_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications/5", nil)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.GetNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestGetNotification_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.GetNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestGetNotification_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return((*models.Notification)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.GetNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestGetNotification_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/invalid/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("invalid", "5")

	err := h.GetNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestGetNotification_InvalidNotificationID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1/notifications/invalid", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "invalid")

	err := h.GetNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
