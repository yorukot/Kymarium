package statuspage

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils/response"
	"go.uber.org/zap"
)

type publicTimelinePoint struct {
	Day     time.Time `json:"day"`
	Success int64     `json:"success"`
	Fail    int64     `json:"fail"`
}

type publicStatusPageMonitor struct {
	ID          string                       `json:"id"`
	MonitorID   string                       `json:"monitor_id"`
	GroupID     *string                      `json:"group_id,omitempty"`
	Name        string                       `json:"name"`
	Type        models.StatusPageElementType `json:"type"`
	SortOrder   int                          `json:"sort_order"`
	Status      string                       `json:"status,omitempty"`
	UptimeSLI30 float64                      `json:"uptime_sli_30,omitempty"`
	UptimeSLI60 float64                      `json:"uptime_sli_60,omitempty"`
	UptimeSLI90 float64                      `json:"uptime_sli_90,omitempty"`
	Timeline    []publicTimelinePoint        `json:"timeline,omitempty"`
}

type publicStatusPageElement struct {
	ID          string                       `json:"id"`
	Name        string                       `json:"name"`
	Type        models.StatusPageElementType `json:"type"`
	SortOrder   int                          `json:"sort_order"`
	Status      string                       `json:"status,omitempty"`
	Monitor     bool                         `json:"monitor"`
	MonitorID   *string                      `json:"monitor_id,omitempty"`
	UptimeSLI30 float64                      `json:"uptime_sli_30,omitempty"`
	UptimeSLI60 float64                      `json:"uptime_sli_60,omitempty"`
	UptimeSLI90 float64                      `json:"uptime_sli_90,omitempty"`
	Timeline    []publicTimelinePoint        `json:"timeline,omitempty"`
	Monitors    []publicStatusPageMonitor    `json:"monitors"`
}

type publicIncidentResponse struct {
	models.Incident
	Timeline  []models.EventTimeline `json:"timeline"`
	MonitorID string                 `json:"monitor_id"`
}

type publicStatusPageResponse struct {
	StatusPage models.StatusPage         `json:"status_page"`
	Elements   []publicStatusPageElement `json:"elements"`
	Incidents  []publicIncidentResponse  `json:"incidents"`
}

// GetPublicStatusPage godoc
// @Summary Get public status page
// @Description Fetches a public status page by slug with computed status/timeline data
// @Tags status-pages
// @Produce json
// @Param slug path string true "Status Page Slug"
// @Success 200 {object} response.SuccessResponse "Public status page returned"
// @Failure 404 {object} response.ErrorResponse "Status page not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /status-pages/{slug} [get]
func (h *Handler) GetPublicStatusPage(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
	}

	tx, err := h.Repo.StartTransaction(c.Request().Context())
	if err != nil {
		zap.L().Error("Failed to begin transaction", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to begin transaction")
	}
	defer h.Repo.DeferRollback(c.Request().Context(), tx)

	page, err := h.Repo.GetStatusPageBySlug(c.Request().Context(), tx, slug)
	if err != nil {
		zap.L().Error("Failed to get status page by slug", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get status page")
	}
	if page == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Status page not found")
	}

	groups, err := h.Repo.ListStatusPageGroupsByStatusPageID(c.Request().Context(), tx, page.ID)
	if err != nil {
		zap.L().Error("Failed to list status page groups", zap.Error(err), zap.Int64("status_page_id", page.ID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list status page groups")
	}

	monitors, err := h.Repo.ListStatusPageMonitorsByStatusPageID(c.Request().Context(), tx, page.ID)
	if err != nil {
		zap.L().Error("Failed to list status page monitors", zap.Error(err), zap.Int64("status_page_id", page.ID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list status page monitors")
	}

	monitorIDs := make([]int64, 0, len(monitors))
	monitorSeen := make(map[int64]struct{}, len(monitors))
	for _, m := range monitors {
		if _, exists := monitorSeen[m.MonitorID]; exists {
			continue
		}
		monitorSeen[m.MonitorID] = struct{}{}
		monitorIDs = append(monitorIDs, m.MonitorID)
	}

	monitorRows, err := h.Repo.ListMonitorsByIDs(c.Request().Context(), tx, page.TeamID, monitorIDs)
	if err != nil {
		zap.L().Error("Failed to list monitors by ids", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list monitors")
	}

	monitorByID := make(map[int64]models.Monitor, len(monitorRows))
	for _, monitor := range monitorRows {
		monitorByID[monitor.ID] = monitor
	}

	incidents, err := h.Repo.ListPublicIncidentsByMonitorIDs(c.Request().Context(), tx, monitorIDs)
	if err != nil {
		zap.L().Error("Failed to list public incidents", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list incidents")
	}

	openPublicIncident := make(map[int64]bool)
	incidentIDs := make([]int64, 0, len(incidents))
	incidentResponses := make([]publicIncidentResponse, 0, len(incidents))
	for _, incident := range incidents {
		if incident.Status != models.IncidentStatusResolved {
			openPublicIncident[incident.MonitorID] = true
		}
		incidentIDs = append(incidentIDs, incident.ID)
		incidentResponses = append(incidentResponses, publicIncidentResponse{
			Incident:  incident.Incident,
			MonitorID: formatID(incident.MonitorID),
		})
	}

	eventTimelines, err := h.Repo.ListPublicEventTimelinesByIncidentIDs(c.Request().Context(), tx, incidentIDs)
	if err != nil {
		zap.L().Error("Failed to list public incident timelines", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list incident timelines")
	}

	start, end := publicTimelineWindow()
	dailySummaries, err := h.Repo.ListMonitorDailySummaryByMonitorIDs(c.Request().Context(), tx, monitorIDs, start, end)
	if err != nil {
		zap.L().Error("Failed to list daily summaries", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list timeline data")
	}

	if err := h.Repo.CommitTransaction(c.Request().Context(), tx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to commit transaction")
	}

	days := buildTimelineDays(start, end)
	publicTimelinesByIncident := make(map[int64][]models.EventTimeline)
	for _, event := range eventTimelines {
		publicTimelinesByIncident[event.IncidentID] = append(publicTimelinesByIncident[event.IncidentID], event)
	}

	for i := range incidentResponses {
		incidentID := incidentResponses[i].Incident.ID
		if timeline, ok := publicTimelinesByIncident[incidentID]; ok {
			incidentResponses[i].Timeline = timeline
		} else {
			incidentResponses[i].Timeline = []models.EventTimeline{}
		}
	}
	perMonitorDaily := buildDailyIndex(dailySummaries)

	groupMonitorIDs := make(map[int64][]int64, len(groups))
	for _, m := range monitors {
		if m.GroupID == nil {
			continue
		}
		groupMonitorIDs[*m.GroupID] = append(groupMonitorIDs[*m.GroupID], m.MonitorID)
	}

	groupMonitorResponses := make(map[int64][]publicStatusPageMonitor, len(groups))
	ungroupedMonitors := make([]publicStatusPageMonitor, 0)
	for _, monitor := range monitors {
		status := computeMonitorStatus(monitor.MonitorID, monitorByID, openPublicIncident)
		timeline, sli30, sli60, sli90 := buildTimelineSummary([]int64{monitor.MonitorID}, days, perMonitorDaily)

		var groupID *string
		if monitor.GroupID != nil {
			id := formatID(*monitor.GroupID)
			groupID = &id
		}

		responseMonitor := publicStatusPageMonitor{
			ID:        formatID(monitor.ID),
			MonitorID: formatID(monitor.MonitorID),
			GroupID:   groupID,
			Name:      monitor.Name,
			Type:      monitor.Type,
			SortOrder: monitor.SortOrder,
			Status:    status,
		}

		if monitor.Type == models.StatusPageElementTypeHistoricalTimeline {
			responseMonitor.Timeline = timeline
			responseMonitor.UptimeSLI30 = sli30
			responseMonitor.UptimeSLI60 = sli60
			responseMonitor.UptimeSLI90 = sli90
		}

		if monitor.GroupID == nil {
			ungroupedMonitors = append(ungroupedMonitors, responseMonitor)
		} else {
			groupMonitorResponses[*monitor.GroupID] = append(groupMonitorResponses[*monitor.GroupID], responseMonitor)
		}
	}

	elements := make([]publicStatusPageElement, 0, len(groups)+len(ungroupedMonitors))
	for _, group := range groups {
		monitorIDs := groupMonitorIDs[group.ID]
		status := computeGroupStatus(monitorIDs, monitorByID, openPublicIncident)
		timeline, sli30, sli60, sli90 := buildTimelineSummary(monitorIDs, days, perMonitorDaily)

		monitorList := groupMonitorResponses[group.ID]
		if monitorList == nil {
			monitorList = []publicStatusPageMonitor{}
		}

		responseElement := publicStatusPageElement{
			ID:        formatID(group.ID),
			Name:      group.Name,
			Type:      group.Type,
			SortOrder: group.SortOrder,
			Status:    status,
			Monitor:   false,
			Monitors:  monitorList,
		}

		if group.Type == models.StatusPageElementTypeHistoricalTimeline {
			responseElement.Timeline = timeline
			responseElement.UptimeSLI30 = sli30
			responseElement.UptimeSLI60 = sli60
			responseElement.UptimeSLI90 = sli90
		}

		elements = append(elements, responseElement)
	}

	for _, monitor := range ungroupedMonitors {
		monitorID := monitor.MonitorID
		element := publicStatusPageElement{
			ID:        monitor.ID,
			Name:      monitor.Name,
			Type:      monitor.Type,
			SortOrder: monitor.SortOrder,
			Status:    monitor.Status,
			Monitor:   true,
			MonitorID: &monitorID,
			Monitors:  []publicStatusPageMonitor{},
		}

		if monitor.Type == models.StatusPageElementTypeHistoricalTimeline {
			element.Timeline = monitor.Timeline
			element.UptimeSLI30 = monitor.UptimeSLI30
			element.UptimeSLI60 = monitor.UptimeSLI60
			element.UptimeSLI90 = monitor.UptimeSLI90
		}

		elements = append(elements, element)
	}

	resp := publicStatusPageResponse{
		StatusPage: *page,
		Elements:   elements,
		Incidents:  incidentResponses,
	}

	return c.JSON(http.StatusOK, response.Success("Status page returned", resp))
}

func publicTimelineWindow() (time.Time, time.Time) {
	now := time.Now().UTC()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
	start := end.AddDate(0, 0, -90)
	return start, end
}

func buildTimelineDays(start time.Time, end time.Time) []time.Time {
	days := make([]time.Time, 0, 90)
	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		days = append(days, day)
	}
	return days
}

type dailyCounts struct {
	total int64
	good  int64
}

func buildDailyIndex(summaries []models.MonitorDailySummary) map[int64]map[time.Time]dailyCounts {
	index := make(map[int64]map[time.Time]dailyCounts)
	for _, summary := range summaries {
		monitorMap, exists := index[summary.MonitorID]
		if !exists {
			monitorMap = make(map[time.Time]dailyCounts)
			index[summary.MonitorID] = monitorMap
		}
		monitorMap[summary.Day.UTC()] = dailyCounts{
			total: summary.TotalCount,
			good:  summary.GoodCount,
		}
	}
	return index
}

func buildTimelineSummary(monitorIDs []int64, days []time.Time, daily map[int64]map[time.Time]dailyCounts) ([]publicTimelinePoint, float64, float64, float64) {
	points := make([]publicTimelinePoint, 0, len(days))
	dayTotals := make([]dailyCounts, 0, len(days))

	for _, day := range days {
		var dayTotal int64
		var dayGood int64
		for _, monitorID := range monitorIDs {
			if byDay, ok := daily[monitorID]; ok {
				if counts, ok := byDay[day]; ok {
					dayTotal += counts.total
					dayGood += counts.good
				}
			}
		}
		dayTotals = append(dayTotals, dailyCounts{
			total: dayTotal,
			good:  dayGood,
		})
		points = append(points, publicTimelinePoint{
			Day:     day,
			Success: dayGood,
			Fail:    dayTotal - dayGood,
		})
	}

	good30, total30 := sumLast(dayTotals, 30)
	good60, total60 := sumLast(dayTotals, 60)
	good90, total90 := sumLast(dayTotals, 90)

	return points, percentage(good30, total30), percentage(good60, total60), percentage(good90, total90)
}

func sumLast(counts []dailyCounts, days int) (int64, int64) {
	if days > len(counts) {
		days = len(counts)
	}
	if days <= 0 {
		return 0, 0
	}
	start := len(counts) - days
	var total int64
	var good int64
	for i := start; i < len(counts); i++ {
		total += counts[i].total
		good += counts[i].good
	}
	return good, total
}

func computeMonitorStatus(monitorID int64, monitorByID map[int64]models.Monitor, openPublicIncident map[int64]bool) string {
	if openPublicIncident[monitorID] {
		return "down"
	}

	monitor, ok := monitorByID[monitorID]
	if !ok {
		return "down"
	}

	if monitor.Status == models.MonitorStatusDown {
		return "down"
	}

	return "up"
}

func computeGroupStatus(monitorIDs []int64, monitorByID map[int64]models.Monitor, openPublicIncident map[int64]bool) string {
	if len(monitorIDs) == 0 {
		return "up"
	}
	for _, monitorID := range monitorIDs {
		if computeMonitorStatus(monitorID, monitorByID, openPublicIncident) == "down" {
			return "down"
		}
	}
	return "up"
}

func formatID(id int64) string {
	return strconv.FormatInt(id, 10)
}

func percentage(good int64, total int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(good) / float64(total)) * 100
}
