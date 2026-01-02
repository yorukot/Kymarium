package team

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/kymarium/internal/testutil"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
)

func TestDeleteTeam_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, int64(7), int64(123)).
		Return(&models.TeamMember{Role: models.MemberRoleOwner}, nil)
	mockRepo.On("DeleteTeam", mock.Anything, mock.Anything, int64(7)).Return(nil)

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodDelete, "/teams/7", nil)
	c.SetParamNames("id")
	c.SetParamValues("7")
	testutil.Authenticate(c, 123)

	err := h.DeleteTeam(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteTeam_Forbidden(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&models.TeamMember{Role: models.MemberRoleMember}, nil)

	h := &Handler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/7", nil)
	c.SetParamNames("id")
	c.SetParamValues("7")
	testutil.Authenticate(c, 555)

	err := h.DeleteTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusForbidden, httpErr.Code)
}

func TestDeleteTeam_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetTeamMemberByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return((*models.TeamMember)(nil), nil)

	h := &Handler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/7", nil)
	c.SetParamNames("id")
	c.SetParamValues("7")
	testutil.Authenticate(c, 555)

	err := h.DeleteTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestDeleteTeam_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &Handler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodDelete, "/teams/7", nil)
	c.SetParamNames("id")
	c.SetParamValues("7")

	err := h.DeleteTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
