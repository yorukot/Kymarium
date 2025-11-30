package notification

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/knocker/internal/testutil"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

func TestUpdateNotification_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "Old Name", Type: "discord"}, nil)
	mockRepo.On("UpdateNotification", mock.Anything, mock.Anything, mock.Anything).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "New Name", Type: "telegram"}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"type":"telegram","name":"New Name","config":{"token":"abc123"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Notification updated successfully", resp["message"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "New Name", data["name"])
}

func TestUpdateNotification_SuccessAsAdmin(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleAdmin}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "Old", Type: "discord"}, nil)
	mockRepo.On("UpdateNotification", mock.Anything, mock.Anything, mock.Anything).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "New", Type: "discord"}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateNotification_PartialUpdate(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "Original", Type: "discord"}, nil)
	mockRepo.On("UpdateNotification", mock.Anything, mock.Anything, mock.Anything).
		Return(&models.Notification{ID: 5, TeamID: 1, Name: "Updated Only Name", Type: "discord"}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"Updated Only Name"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateNotification_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestUpdateNotification_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestUpdateNotification_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestUpdateNotification_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return((*models.Notification)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestUpdateNotification_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/invalid/notifications/5",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("invalid", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestUpdateNotification_InvalidNotificationID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/invalid",
		strings.NewReader(`{"name":"New"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "invalid")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestUpdateNotification_NoFieldsProvided(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPatch, "/teams/1/notifications/5",
		strings.NewReader(`{}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.UpdateNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
