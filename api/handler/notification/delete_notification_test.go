package notification

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/knocker/internal/testutil"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

func TestDeleteNotification_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("DeleteNotification", mock.Anything, mock.Anything, int64(1), int64(5)).Return(nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.DeleteNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteNotification_SuccessAsAdmin(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleAdmin}, nil)
	mockRepo.On("DeleteNotification", mock.Anything, mock.Anything, int64(1), int64(5)).Return(nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.DeleteNotification(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteNotification_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/5", nil)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.DeleteNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestDeleteNotification_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.DeleteNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestDeleteNotification_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.DeleteNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestDeleteNotification_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/invalid/notifications/5", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("invalid", "5")

	err := h.DeleteNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestDeleteNotification_InvalidNotificationID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/1/notifications/invalid", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "invalid")

	err := h.DeleteNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
