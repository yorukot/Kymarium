package monitor

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/utils"
)

type monitorResponse struct {
	ID                string               `json:"id"`
	TeamID            string               `json:"team_id"`
	Name              string               `json:"name"`
	Type              models.MonitorType   `json:"type"`
	Config            json.RawMessage      `json:"config"`
	Interval          int                  `json:"interval"`
	Status            models.MonitorStatus `json:"status"`
	UptimeSLI30       *float64             `json:"uptime_sli_30,omitempty"`
	LastChecked       time.Time            `json:"last_checked"`
	NextCheck         time.Time            `json:"next_check"`
	FailureThreshold  int16                `json:"failure_threshold"`
	RecoveryThreshold int16                `json:"recovery_threshold"`
	RegionIDs         []string             `json:"regions"`
	NotificationIDs   []string             `json:"notification"`
	Incidents         []incidentResponse   `json:"incidents,omitempty"`
	UpdatedAt         time.Time            `json:"updated_at"`
	CreatedAt         time.Time            `json:"created_at"`
}

type incidentResponse struct {
	ID         string                `json:"id"`
	MonitorID  string                `json:"monitor_id"`
	Title      *string               `json:"title,omitempty"`
	Status     models.IncidentStatus `json:"status"`
	StartedAt  time.Time             `json:"started_at"`
	ResolvedAt *time.Time            `json:"resolved_at,omitempty"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type notificationIDList = utils.IDList
type regionIDList = utils.IDList

func newMonitorResponse(m models.Monitor) monitorResponse {
	return newMonitorResponseWithUptime(m, nil)
}

func newMonitorResponseWithUptime(m models.Monitor, uptimeSLI30 *float64) monitorResponse {
	return monitorResponse{
		ID:                strconv.FormatInt(m.ID, 10),
		TeamID:            strconv.FormatInt(m.TeamID, 10),
		Name:              m.Name,
		Type:              m.Type,
		Config:            m.Config,
		Interval:          m.Interval,
		Status:            m.Status,
		UptimeSLI30:       uptimeSLI30,
		LastChecked:       m.LastChecked,
		NextCheck:         m.NextCheck,
		FailureThreshold:  m.FailureThreshold,
		RecoveryThreshold: m.RecoveryThreshold,
		RegionIDs:         formatRegionIDs(m.RegionIDs),
		NotificationIDs:   formatNotificationIDs(m.NotificationIDs),
		Incidents:         []incidentResponse{},
		UpdatedAt:         m.UpdatedAt,
		CreatedAt:         m.CreatedAt,
	}
}

func newMonitorResponseWithIncidents(m models.MonitorWithIncidents, uptimeSLI30 *float64) monitorResponse {
	resp := newMonitorResponseWithUptime(m.Monitor, uptimeSLI30)
	resp.Incidents = formatIncidents(m.ID, m.Incidents)
	return resp
}

func newMonitorResponses(monitors []models.Monitor) []monitorResponse {
	responses := make([]monitorResponse, len(monitors))
	for i, monitor := range monitors {
		responses[i] = newMonitorResponse(monitor)
	}
	return responses
}

func newMonitorResponsesWithIncidents(monitors []models.MonitorWithIncidents, uptimeByMonitor map[int64]*float64) []monitorResponse {
	responses := make([]monitorResponse, len(monitors))
	for i, monitor := range monitors {
		var uptime *float64
		if uptimeByMonitor != nil {
			uptime = uptimeByMonitor[monitor.ID]
		}
		responses[i] = newMonitorResponseWithIncidents(monitor, uptime)
	}
	return responses
}

func formatNotificationIDs(ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result
}

func formatRegionIDs(ids []int64) []string {
	if len(ids) == 0 {
		return []string{}
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = strconv.FormatInt(id, 10)
	}
	return result
}

func formatIncidents(monitorID int64, incidents []models.Incident) []incidentResponse {
	if len(incidents) == 0 {
		return []incidentResponse{}
	}

	result := make([]incidentResponse, len(incidents))
	for i, incident := range incidents {
		result[i] = incidentResponse{
			ID:         strconv.FormatInt(incident.ID, 10),
			MonitorID:  strconv.FormatInt(monitorID, 10),
			Title:      incident.Title,
			Status:     incident.Status,
			StartedAt:  incident.StartedAt,
			ResolvedAt: incident.ResolvedAt,
			CreatedAt:  incident.CreatedAt,
			UpdatedAt:  incident.UpdatedAt,
		}
	}
	return result
}
