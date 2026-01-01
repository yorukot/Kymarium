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

func TestCreateTeam_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	var capturedTeam models.Team
	var capturedMember models.TeamMember
	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("CreateTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedTeam = args.Get(2).(models.Team)
	})
	mockRepo.On("CreateTeamMember", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedMember = args.Get(2).(models.TeamMember)
	})

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodPost, "/teams", strings.NewReader(`{"name":"Acme"}`))
	testutil.SetJSONHeader(c)
	testutil.Authenticate(c, 123)

	err := h.CreateTeam(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, capturedTeam.ID, capturedMember.TeamID, "member should reference created team")
	require.Equal(t, int64(123), capturedMember.UserID)
	require.Equal(t, models.MemberRoleOwner, capturedMember.Role)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Team created successfully", resp["message"])
}

func TestCreateTeam_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &Handler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodPost, "/teams", strings.NewReader(`{"name":"Acme"}`))
	testutil.SetJSONHeader(c)

	err := h.CreateTeam(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
