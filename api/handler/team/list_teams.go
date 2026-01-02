package team

import (
	"net/http"

	"github.com/labstack/echo/v4"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
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
func (h *Handler) ListTeams(c echo.Context) error {
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

	teams, err := h.Repo.ListTeamsByUserID(c.Request().Context(), tx, *userID)
	if err != nil {
		zap.L().Error("Failed to list teams", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list teams")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Teams retrieved successfully", teams))
}
