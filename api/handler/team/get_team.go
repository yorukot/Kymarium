package team

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/repository"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// GetTeam godoc
// @Summary Get a team
// @Description Retrieves a single team the authenticated user belongs to
// @Tags teams
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} response.SuccessResponse "Team retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{id} [get]
func (h *TeamHandler) GetTeam(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

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

	team, err := repository.GetTeamForUser(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team")
	}

	if team == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	if err := repository.CommitTransaction(tx, c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Team retrieved successfully", team))
}
