package monitor

import (
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	authutil "github.com/yorukot/kymarium/utils/auth"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

const defaultAnalyticsLookback = 24 * time.Hour

type analyticsWindow struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Bucket string    `json:"bucket"`
}

type analyticsSummary struct {
	TotalCount int64   `json:"total_count"`
	GoodCount  int64   `json:"good_count"`
	UptimePct  float64 `json:"uptime_pct"`
	P50Ms      float64 `json:"p50_ms"`
	P75Ms      float64 `json:"p75_ms"`
	P90Ms      float64 `json:"p90_ms"`
	P95Ms      float64 `json:"p95_ms"`
	P99Ms      float64 `json:"p99_ms"`
}

type analyticsRegionSummary struct {
	RegionID   string  `json:"region_id"`
	TotalCount int64   `json:"total_count"`
	GoodCount  int64   `json:"good_count"`
	UptimePct  float64 `json:"uptime_pct"`
	P50Ms      float64 `json:"p50_ms"`
	P75Ms      float64 `json:"p75_ms"`
	P90Ms      float64 `json:"p90_ms"`
	P95Ms      float64 `json:"p95_ms"`
	P99Ms      float64 `json:"p99_ms"`
}

type analyticsSeriesPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	RegionID   string    `json:"region_id"`
	TotalCount int64     `json:"total_count"`
	GoodCount  int64     `json:"good_count"`
	UptimePct  float64   `json:"uptime_pct"`
	P50Ms      float64   `json:"p50_ms"`
	P75Ms      float64   `json:"p75_ms"`
	P90Ms      float64   `json:"p90_ms"`
	P95Ms      float64   `json:"p95_ms"`
	P99Ms      float64   `json:"p99_ms"`
}

type monitorAnalyticsResponse struct {
	Monitor   monitorResponse          `json:"monitor"`
	Window    analyticsWindow          `json:"window"`
	Summary   analyticsSummary         `json:"summary"`
	Regions   []analyticsRegionSummary `json:"regions"`
	Series    []analyticsSeriesPoint   `json:"series"`
	Incidents []incidentResponse       `json:"incidents"`
}

// GetAnalytics godoc
// @Summary Get monitor analytics
// @Description Returns uptime and latency analytics for a monitor within a window (default last 24h, bucket 30m)
// @Tags monitors
// @Produce json
// @Param teamID path string true "Team ID"
// @Param id path string true "Monitor ID"
// @Param start query string false "Start time (ISO8601)"
// @Param end query string false "End time (ISO8601)"
// @Param bucket query string false "Bucket duration, only 30m supported"
// @Param region_id query string false "Region ID to filter"
// @Success 200 {object} response.SuccessResponse "Analytics returned"
// @Failure 400 {object} response.ErrorResponse "Invalid parameters"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Monitor or team not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /teams/{teamID}/monitors/{id}/analytics [get]
func (h *Handler) GetAnalytics(c echo.Context) error {
	teamID, err := strconv.ParseInt(c.Param("teamID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid team ID")
	}

	monitorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid monitor ID")
	}

	userID, err := authutil.GetUserIDFromContext(c)
	if err != nil {
		zap.L().Error("Failed to parse user ID from context", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	if userID == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	end := time.Now().UTC()
	if endStr := c.QueryParam("end"); endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid end time format")
		}
	}

	start := end.Add(-defaultAnalyticsLookback)
	if startStr := c.QueryParam("start"); startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid start time format")
		}
	}

	if !start.Before(end) {
		return echo.NewHTTPError(http.StatusBadRequest, "start must be before end")
	}

	bucketParam := c.QueryParam("bucket")
	if bucketParam == "" {
		bucketParam = "30m"
	}
	if bucketParam != "30m" {
		return echo.NewHTTPError(http.StatusBadRequest, "Only bucket=30m is supported")
	}

	var regionFilter *int64
	if regionParam := c.QueryParam("region_id"); regionParam != "" {
		regionVal, parseErr := strconv.ParseInt(regionParam, 10, 64)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid region_id")
		}
		regionFilter = &regionVal
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

	monitor, err := h.Repo.GetMonitorByID(c.Request().Context(), tx, teamID, monitorID)
	if err != nil {
		zap.L().Error("Failed to get monitor", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get monitor")
	}
	if monitor == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Monitor not found")
	}

	if regionFilter != nil && !slices.Contains(monitor.RegionIDs, *regionFilter) {
		return echo.NewHTTPError(http.StatusBadRequest, "region_id is not associated with this monitor")
	}

	buckets, err := h.Repo.GetMonitorAnalytics(c.Request().Context(), tx, monitorID, start, end, regionFilter)
	if err != nil {
		zap.L().Error("Failed to fetch analytics", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch analytics")
	}

	incidents, err := h.Repo.ListIncidentsByMonitorIDWithinRange(c.Request().Context(), tx, monitorID, start, end)
	if err != nil {
		zap.L().Error("Failed to list incidents", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list incidents")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	resp := buildAnalyticsResponse(*monitor, buckets, incidents, start, end, bucketParam)

	return c.JSON(http.StatusOK, response.Success("Analytics returned", resp))
}

func buildAnalyticsResponse(monitor models.Monitor, buckets []models.MonitorAnalyticsBucket, incidents []models.Incident, start, end time.Time, bucket string) monitorAnalyticsResponse {
	overall := analyticsSummary{}
	regionMap := make(map[int64]*analyticsSummary)
	series := make([]analyticsSeriesPoint, 0, len(buckets))

	for _, bucketRow := range buckets {
		overall.TotalCount += bucketRow.TotalCount
		overall.GoodCount += bucketRow.GoodCount
		weight := float64(bucketRow.TotalCount)
		if weight > 0 {
			overall.P50Ms += bucketRow.P50Ms * weight
			overall.P75Ms += bucketRow.P75Ms * weight
			overall.P90Ms += bucketRow.P90Ms * weight
			overall.P95Ms += bucketRow.P95Ms * weight
			overall.P99Ms += bucketRow.P99Ms * weight
		}

		regSummary, exists := regionMap[bucketRow.RegionID]
		if !exists {
			regSummary = &analyticsSummary{}
			regionMap[bucketRow.RegionID] = regSummary
		}

		regSummary.TotalCount += bucketRow.TotalCount
		regSummary.GoodCount += bucketRow.GoodCount
		if bucketRow.TotalCount > 0 {
			w := float64(bucketRow.TotalCount)
			regSummary.P50Ms += bucketRow.P50Ms * w
			regSummary.P75Ms += bucketRow.P75Ms * w
			regSummary.P90Ms += bucketRow.P90Ms * w
			regSummary.P95Ms += bucketRow.P95Ms * w
			regSummary.P99Ms += bucketRow.P99Ms * w
		}

		point := analyticsSeriesPoint{
			Timestamp:  bucketRow.Bucket,
			RegionID:   strconv.FormatInt(bucketRow.RegionID, 10),
			TotalCount: bucketRow.TotalCount,
			GoodCount:  bucketRow.GoodCount,
			UptimePct:  percentage(bucketRow.GoodCount, bucketRow.TotalCount),
			P50Ms:      bucketRow.P50Ms,
			P75Ms:      bucketRow.P75Ms,
			P90Ms:      bucketRow.P90Ms,
			P95Ms:      bucketRow.P95Ms,
			P99Ms:      bucketRow.P99Ms,
		}

		series = append(series, point)
	}

	if overall.TotalCount > 0 {
		w := float64(overall.TotalCount)
		overall.UptimePct = percentage(overall.GoodCount, overall.TotalCount)
		overall.P50Ms = overall.P50Ms / w
		overall.P75Ms = overall.P75Ms / w
		overall.P90Ms = overall.P90Ms / w
		overall.P95Ms = overall.P95Ms / w
		overall.P99Ms = overall.P99Ms / w
	}

	regions := make([]analyticsRegionSummary, 0, len(regionMap))
	for regionID, summary := range regionMap {
		if summary.TotalCount > 0 {
			w := float64(summary.TotalCount)
			summary.UptimePct = percentage(summary.GoodCount, summary.TotalCount)
			summary.P50Ms = summary.P50Ms / w
			summary.P75Ms = summary.P75Ms / w
			summary.P90Ms = summary.P90Ms / w
			summary.P95Ms = summary.P95Ms / w
			summary.P99Ms = summary.P99Ms / w
		}
		regions = append(regions, analyticsRegionSummary{
			RegionID:   strconv.FormatInt(regionID, 10),
			TotalCount: summary.TotalCount,
			GoodCount:  summary.GoodCount,
			UptimePct:  summary.UptimePct,
			P50Ms:      summary.P50Ms,
			P75Ms:      summary.P75Ms,
			P90Ms:      summary.P90Ms,
			P95Ms:      summary.P95Ms,
			P99Ms:      summary.P99Ms,
		})
	}

	return monitorAnalyticsResponse{
		Monitor: newMonitorResponse(monitor),
		Window: analyticsWindow{
			Start:  start,
			End:    end,
			Bucket: bucket,
		},
		Summary:   overall,
		Regions:   regions,
		Series:    series,
		Incidents: formatIncidents(monitor.ID, incidents),
	}
}

func percentage(good, total int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(good) / float64(total)) * 100
}
