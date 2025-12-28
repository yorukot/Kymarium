package auth

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/utils/config"
	"go.uber.org/zap"
)

// VerifyEmail godoc
// @Summary Verify email
// @Description Verifies a user's email address using the verification token
// @Tags auth
// @Produce json
// @Param token query string true "Verification token"
// @Success 302 {object} response.SuccessResponse "Email verified successfully"
// @Failure 302 {object} response.ErrorResponse "Invalid or expired verification token"
// @Failure 302 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/verify [get]
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	verificationToken := c.QueryParam("token")
	if verificationToken == "" {
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("expired", ""))
	}

	valid, claims, err := validateEmailVerificationToken(verificationToken)
	if err != nil {
		zap.L().Error("Failed to validate verification token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate verification token")
	}
	if !valid {
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("expired", ""))
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("expired", ""))
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(tx, c.Request().Context())

	user, err := h.Repo.GetUserByID(c.Request().Context(), tx, userID)
	if err != nil {
		zap.L().Error("Failed to get user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}
	if user == nil {
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("expired", claims.Email))
	}

	if user.Verified {
		if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
			zap.L().Error("Failed to commit transaction", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
		}
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("success", claims.Email))
	}

	if user.VerifyCode == nil || *user.VerifyCode != verificationToken {
		return c.Redirect(http.StatusFound, buildFrontendVerificationURL("expired", claims.Email))
	}

	if err := h.Repo.UpdateUserVerification(c.Request().Context(), tx, user.ID, true, nil, time.Now()); err != nil {
		zap.L().Error("Failed to update user verification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user verification")
	}

	if err := h.Repo.CommitTransaction(tx, c.Request().Context()); err != nil {
		zap.L().Error("Failed to commit transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.Redirect(http.StatusFound, buildFrontendVerificationURL("success", claims.Email))
}

func buildFrontendVerificationURL(result string, email string) string {
	base := strings.TrimSpace(config.Env().FrontendDomain)
	if base == "" {
		return "/"
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}

	base = strings.TrimRight(base, "/")
	path := "/auth/verify/" + result

	if email == "" {
		return base + path
	}

	return base + path + "?email=" + url.QueryEscape(email)
}
