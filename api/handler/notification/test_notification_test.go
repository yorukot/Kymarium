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

// Note: TestNotification handler tests for success cases are not included because the handler
// calls notificationcore.Send() which is an external service call that cannot be mocked at this level.
// The handler is tested for authorization and validation logic instead.

func TestTestNotification_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications/5/test", nil)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestTestNotification_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications/5/test", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestTestNotification_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications/5/test", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestTestNotification_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("GetNotificationByID", mock.Anything, mock.Anything, int64(1), int64(5)).
		Return((*models.Notification)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications/5/test", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "5")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestTestNotification_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/invalid/notifications/5/test", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("invalid", "5")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestTestNotification_InvalidNotificationID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications/invalid/test", nil)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID", "id")
	c.SetParamValues("1", "invalid")

	err := h.TestNotification(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
