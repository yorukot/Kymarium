package team

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

func TestUpdateTeam_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(10), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleAdmin}, nil)
	mockRepo.On("UpdateTeamName", mock.Anything, mock.Anything, int64(10), "NewName", mock.AnythingOfType("time.Time")).
		Return(&models.Team{ID: 10, Name: "NewName"}, nil)

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPut, "/teams/10", strings.NewReader(`{"name":"NewName"}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("id")
	c.SetParamValues("10")
	testutil.Authenticate(c, 123)

	err := h.UpdateTeam(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Team updated successfully", resp["message"])
}

func TestUpdateTeam_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&models.TeamMember{Role: models.MemberRoleViewer}, nil)

	h := &Handler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPut, "/teams/11", strings.NewReader(`{"name":"NewName"}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("id")
	c.SetParamValues("11")
	testutil.Authenticate(c, 999)

	err := h.UpdateTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestUpdateTeam_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return((*models.TeamMember)(nil), nil)

	h := &Handler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodPut, "/teams/11", strings.NewReader(`{"name":"NewName"}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("id")
	c.SetParamValues("11")
	testutil.Authenticate(c, 999)

	err := h.UpdateTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestUpdateTeam_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &Handler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPut, "/teams/1", strings.NewReader(`{"name":"NewName"}`))
	testutil.SetJSONHeader(c)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.UpdateTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
