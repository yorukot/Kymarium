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

func TestNew_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	var capturedNotification models.Notification
	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			capturedNotification = args.Get(2).(models.Notification)
		})
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"discord","name":"Test Notification","config":{"webhook":"https://example.com"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, int64(1), capturedNotification.TeamID)
	require.Equal(t, models.NotificationType("discord"), capturedNotification.Type)
	require.Equal(t, "Test Notification", capturedNotification.Name)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Notification created successfully", resp["message"])
}

func TestNew_SuccessAsAdmin(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleAdmin}, nil)
	mockRepo.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"telegram","name":"Test Notification","config":{"token":"abc123"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestNew_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"discord","name":"Test","config":{"webhook":"test"}}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestNew_NotMember(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return((*models.TeamMember)(nil), nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"discord","name":"Test","config":{"webhook":"test"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestNew_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(1), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)

	h := &NotificationHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"discord","name":"Test","config":{"webhook":"test"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestNew_InvalidTeamID(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/invalid/notifications",
		strings.NewReader(`{"type":"discord","name":"Test","config":{"webhook":"test"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("invalid")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestNew_InvalidJSON(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`invalid json`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestNew_MissingRequiredFields(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"discord"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestNew_InvalidType(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &NotificationHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams/1/notifications",
		strings.NewReader(`{"type":"invalid","name":"Test","config":{"webhook":"test"}}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)
	c.SetParamNames("teamID")
	c.SetParamValues("1")

	err := h.New(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusBadRequest, httpErr.Code)
}
