package team

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/repository"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListTeams godoc
// @Summary List teams
// @Description Lists teams the authenticated user is a member of
// @Tags teams
// @Produce json
// @Success 200 {object} response.SuccessResponse "Teams retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams [get]
func (h *TeamHandler) ListTeams(c echo.Context) error {
	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	tx, err := repository.StartTransaction(h.DB, c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer repository.DeferRollback(tx, c.Request().Context())

	teams, err := repository.ListTeamsByUserID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to list teams", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list teams")
	}

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Teams retrieved successfully", teams))
}
