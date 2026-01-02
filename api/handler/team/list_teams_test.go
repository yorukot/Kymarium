package team

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/kymarium/internal/testutil"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
)

func TestListTeams_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("ListTeamsByUserID", mock.Anything, mock.Anything, int64(123)).Return([]models.TeamWithRole{
		{Team: models.Team{ID: 1, Name: "One"}, Role: models.MemberRoleOwner},
	}, nil)

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/teams", nil)
	testutil.Authenticate(c, 123)

	err := h.ListTeams(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Teams retrieved successfully", resp["message"])
	data, ok := resp["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 1)
}

func TestListTeams_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &Handler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/teams", nil)

	err := h.ListTeams(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
