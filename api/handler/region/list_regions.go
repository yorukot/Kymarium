package region

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListRegions godoc
// @Summary List regions
// @Description Lists all available monitoring regions
// @Tags regions
// @Produce json
// @Success 200 {object} response.SuccessResponse "Regions retrieved successfully"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /regions [get]
func (h *Handler) ListRegions(c echo.Context) error {
	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	regions, err := h.Repo.ListAllRegions(c.Request().Context(), tx)
	if err != nil {
		zap.L().Error("Failed to list regions", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list regions")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Regions retrieved successfully", regions))
}
