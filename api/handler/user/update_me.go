package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type updateUserRequest struct {
	DisplayName *string `json:"display_name" validate:"omitempty,min=1,max=255"`
	Avatar      *string `json:"avatar" validate:"omitempty,max=2048"`
}

// UpdateMe godoc
// @Summary Update current user
// @Description Updates the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Param request body updateUserRequest true "User update request"
// @Success 200 {object} response.SuccessResponse "User updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me [patch]
func (h *Handler) UpdateMe(c echo.Context) error {
	var req updateUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.DisplayName == nil && req.Avatar == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No fields to update")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	user, err := h.Repo.GetUserByID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to get user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	displayName := user.DisplayName
	if req.DisplayName != nil {
		trimmed := strings.TrimSpace(*req.DisplayName)
		if trimmed == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Display name is required")
		}
		displayName = trimmed
	}

	avatar := user.Avatar
	if req.Avatar != nil {
		trimmed := strings.TrimSpace(*req.Avatar)
		if trimmed == "" {
			avatar = nil
		} else {
			avatar = &trimmed
		}
	}

	updatedUser, err := h.Repo.UpdateUserProfile(c.Request().Context(), tx, *userID, displayName, avatar, time.Now())
	if err != nil {
		zap.L().Error("Failed to update user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	if updatedUser == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("User updated successfully", newUserResponse(updatedUser)))
}
