package team

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

func TestGetTeam_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetTeamForUser", mock.Anything, mock.Anything, int64(5), int64(123)).
		Return(&models.TeamWithRole{Team: models.Team{ID: 5, Name: "Five"}, Role: models.MemberRoleMember}, nil)

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/teams/5", nil)
	c.SetParamNames("id")
	c.SetParamValues("5")
	testutil.Authenticate(c, 123)

	err := h.GetTeam(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Team retrieved successfully", resp["message"])
	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "Five", data["name"])
}

func TestGetTeam_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamForUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return((*models.TeamWithRole)(nil), nil)

	h := &Handler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/9", nil)
	c.SetParamNames("id")
	c.SetParamValues("9")
	testutil.Authenticate(c, 321)

	err := h.GetTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestGetTeam_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &Handler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.GetTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
