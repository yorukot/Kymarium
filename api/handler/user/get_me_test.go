package user

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/knocker/internal/testutil"
	"github.com/yorukot/knocker/models"
	"github.com/yorukot/knocker/repository"
)

func TestGetMe_Success(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	avatar := "https://example.com/avatar.png"
	now := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything, int64(123)).
		Return(&models.User{
			ID:          123,
			DisplayName: "Jane Doe",
			Avatar:      &avatar,
			Verified:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}, nil)

	h := &UserHandler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/users/me", nil)
	testutil.Authenticate(c, 123)

	err := h.GetMe(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "User retrieved successfully", resp["message"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "123", data["id"])
	require.Equal(t, "Jane Doe", data["display_name"])
	require.Equal(t, "https://example.com/avatar.png", data["avatar"])
	require.Equal(t, true, data["verified"])
	require.Contains(t, data, "created_at")
	require.Contains(t, data, "updated_at")
	require.NotContains(t, data, "password_hash")
}

func TestGetMe_NotFound(t *testing.T) {
	testutil.InitTestEnv(t)

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetUserByID", mock.Anything, mock.Anything, int64(456)).
		Return((*models.User)(nil), nil)

	h := &UserHandler{Repo: mockRepo}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/users/me", nil)
	testutil.Authenticate(c, 456)

	err := h.GetMe(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, httpErr.Code)
}

func TestGetMe_Unauthorized(t *testing.T) {
	testutil.InitTestEnv(t)

	h := &UserHandler{Repo: &repository.MockRepository{}}
	c, _ := testutil.NewEchoContext(http.MethodGet, "/users/me", nil)

	err := h.GetMe(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
}
