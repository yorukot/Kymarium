package auth

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

type authStatusResponse struct {
	UserID string `json:"user_id"`
}

// Status godoc
// @Summary Check authentication status
// @Description Returns 200 when the user is authenticated via access token cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse "Authenticated"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Router /auth/status [get]
func (h *Handler) Status(c echo.Context) error {
	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	return c.JSON(http.StatusOK, response.Success("Authenticated", authStatusResponse{
		UserID: strconv.FormatInt(*userID, 10),
	}))
}
