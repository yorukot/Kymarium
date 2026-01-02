package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/encrypt"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type updatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=255"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=255"`
}

// UpdatePassword godoc
// @Summary Update password
// @Description Updates the authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Param request body updatePasswordRequest true "Password update request"
// @Success 200 {object} response.SuccessResponse "Password updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/password [patch]
func (h *Handler) UpdatePassword(c echo.Context) error {
	var req updatePasswordRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if strings.TrimSpace(req.CurrentPassword) == "" || strings.TrimSpace(req.NewPassword) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
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

	if user.PasswordHash == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Password is not set for this account")
	}

	match, err := encrypt.ComparePasswordAndHash(req.CurrentPassword, *user.PasswordHash)
	if err != nil {
		zap.L().Error("Failed to compare password and hash", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to compare password")
	}

	if !match {
		return echo.NewHTTPError(http.StatusBadRequest, "Current password is incorrect")
	}

	if req.CurrentPassword == req.NewPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password must be different")
	}

	newHash, err := encrypt.CreateArgon2idHash(req.NewPassword)
	if err != nil {
		zap.L().Error("Failed to hash password", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	if err := h.Repo.UpdateUserPasswordHash(c.Request().Context(), tx, *userID, newHash, time.Now()); err != nil {
		zap.L().Error("Failed to update password hash", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Password updated successfully"))
}
