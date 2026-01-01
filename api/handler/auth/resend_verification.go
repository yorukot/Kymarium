package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/utils/config"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type resendVerificationRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

// ResendVerification godoc
// @Summary Resend verification email
// @Description Resends the email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param request body resendVerificationRequest true "Resend verification request"
// @Success 200 {object} response.SuccessResponse "Verification email sent"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/verify/resend [post]
func (h *Handler) ResendVerification(c echo.Context) error {
	var req resendVerificationRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validator.New().Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	user, err := h.Repo.GetUserByEmail(c.Request().Context(), tx, req.Email)
	if err != nil {
		zap.L().Error("Failed to get user by email", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user by email")
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	if user.Verified {
		if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
			zap.L().Error("Failed to commit transaction", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}
		return c.JSON(http.StatusOK, response.SuccessMessage("Email already verified"))
	}

	if !config.Env().SMTPEnabled {
		if err := h.Repo.UpdateUserVerification(c.Request().Context(), tx, user.ID, true, nil, time.Now()); err != nil {
			zap.L().Error("Failed to auto-verify user", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify user")
		}

		if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
			zap.L().Error("Failed to commit transaction", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}

		return c.JSON(http.StatusOK, response.SuccessMessage("Email verified"))
	}

	verifyToken, err := generateEmailVerificationToken(user.ID, req.Email)
	if err != nil {
		zap.L().Error("Failed to generate verification token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate verification token")
	}

	if err := h.Repo.UpdateUserVerification(c.Request().Context(), tx, user.ID, false, &verifyToken, time.Now()); err != nil {
		zap.L().Error("Failed to update verification token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update verification token")
	}

	if err := sendVerificationEmail(req.Email, verifyToken); err != nil {
		zap.L().Error("Failed to send verification email", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send verification email")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.SuccessMessage("Verification email sent"))
}
