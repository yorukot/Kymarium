package monitor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/knocker/models"
	authutil "github.com/yorukot/knocker/utils/auth"
	"github.com/yorukot/knocker/utils/response"
	"go.uber.org/zap"
)

// ListMonitors godoc
// @Summary List monitors
// @Description Lists monitors for a team the user belongs to
// @Tags monitors
// @Produce json
// @Param teamID path string true "Team ID"
// @Success 200 {object} response.SuccessResponse "Monitors retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid team ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors [get]
func (h *Handler) ListMonitors(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	const incidentLimit = 5

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

	member, err := h.Repo.GetTeamMemberByUserID(c.Request().Context(), tx, teamID, *userID)
	if err != nil {
		zap.L().Error("Failed to get team membership", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get team membership")
	}

	if member == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Team not found")
	}

	monitors, err := h.Repo.ListMonitorsByTeamID(c.Request().Context(), tx, teamID)
	if err != nil {
		zap.L().Error("Failed to list monitors", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list monitors")
	}

	monitorIDs := make([]int64, 0, len(monitors))
	for _, monitor := range monitors {
		monitorIDs = append(monitorIDs, monitor.ID)
	}

	start, end := last30DayWindowUTC()
	dailySummaries, err := h.Repo.ListMonitorDailySummaryByMonitorIDs(c.Request().Context(), tx, monitorIDs, start, end)
	if err != nil {
		zap.L().Error("Failed to list daily summaries", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list daily summaries")
	}
	uptimeByMonitor := buildUptimeSLI30(dailySummaries)

	monitorsWithIncidents := make([]models.MonitorWithIncidents, 0, len(monitors))
	for _, monitor := range monitors {
		incidents, err := h.Repo.ListIncidentsByMonitorID(c.Request().Context(), tx, monitor.ID)
		if err != nil {
			zap.L().Error("Failed to list incidents for monitor", zap.Error(err), zap.Int64("monitor_id", monitor.ID))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list incidents for monitor")
		}

		if len(incidents) > incidentLimit {
			incidents = incidents[:incidentLimit]
		}

		monitorsWithIncidents = append(monitorsWithIncidents, models.MonitorWithIncidents{
			Monitor:   monitor,
			Incidents: incidents,
		})
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	return c.JSON(http.StatusOK, response.Success("Monitors retrieved successfully", newMonitorResponsesWithIncidents(monitorsWithIncidents, uptimeByMonitor)))
}

func last30DayWindowUTC() (time.Time, time.Time) {
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start := end.AddDate(0, 0, -30)
	return start, end
}

type uptimeTotals struct {
	good  int64
	total int64
}

func buildUptimeSLI30(summaries []models.MonitorDailySummary) map[int64]*float64 {
	totals := make(map[int64]uptimeTotals)
	for _, summary := range summaries {
		current := totals[summary.MonitorID]
		current.total += summary.TotalCount
		current.good += summary.GoodCount
		totals[summary.MonitorID] = current
	}

	result := make(map[int64]*float64, len(totals))
	for monitorID, counts := range totals {
		if counts.total == 0 {
			continue
		}
		uptime := (float64(counts.good) / float64(counts.total)) * 100
		result[monitorID] = &uptime
	}

	return result
}
