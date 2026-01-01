package user

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type sessionResponse struct {
	ID        int64     `json:"id,string" example:"175928847299117063"`
	UserAgent *string   `json:"user_agent,omitempty" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"`
	IP        string    `json:"ip,omitempty" example:"192.168.1.100"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	Current   bool      `json:"current" example:"true"`
}

// ListSessions godoc
// @Summary List active sessions
// @Description Lists active refresh token sessions for the authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} response.SuccessResponse "Sessions retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /users/me/sessions [get]
func (h *Handler) ListSessions(c echo.Context) error {
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

	var currentTokenID *int64
	refreshCookie, err := c.Cookie(models.CookieNameRefreshToken)
	if err == nil && refreshCookie != nil && refreshCookie.Value != "" {
		refreshToken, err := h.Repo.GetRefreshTokenByToken(c.Request().Context(), tx, refreshCookie.Value)
		if err != nil {
			zap.L().Error("Failed to get refresh token", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get sessions")
		}
		if refreshToken != nil {
			currentTokenID = &refreshToken.ID
		}
	}

	tokens, err := h.Repo.ListActiveRefreshTokensByUserID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to list refresh tokens", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get sessions")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	sessions := make([]sessionResponse, 0, len(tokens))
	for _, token := range tokens {
		ip := ""
		if token.IP != nil {
			ip = token.IP.String()
		}
		sessions = append(sessions, sessionResponse{
			ID:        token.ID,
			UserAgent: token.UserAgent,
			IP:        ip,
			CreatedAt: token.CreatedAt,
			Current:   currentTokenID != nil && token.ID == *currentTokenID,
		})
	}

	return c.JSON(http.StatusOK, response.Success("Sessions retrieved successfully", sessions))
}
