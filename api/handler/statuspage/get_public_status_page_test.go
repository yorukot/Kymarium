package statuspage

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/kymarium/internal/testutil"
	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
)

func TestGetPublicStatusPage_DeduplicatesIncidents(t *testing.T) {
	testutil.InitTestEnv(t)

	now := time.Date(2024, time.January, 2, 15, 4, 5, 0, time.UTC)

	page := &models.StatusPage{
		ID:        1,
		TeamID:    10,
		Title:     "Demo",
		Slug:      "demo",
		CreatedAt: now,
		UpdatedAt: now,
	}

	monitors := []models.StatusPageMonitor{
		{ID: 1, StatusPageID: page.ID, MonitorID: 11, Name: "API", Type: models.StatusPageElementTypeCurrentStatusIndicator, SortOrder: 1},
		{ID: 2, StatusPageID: page.ID, MonitorID: 22, Name: "DB", Type: models.StatusPageElementTypeCurrentStatusIndicator, SortOrder: 2},
	}

	monitorRows := []models.Monitor{
		{ID: 11, TeamID: page.TeamID, Name: "API", Type: models.MonitorTypeHTTP, Status: models.MonitorStatusUp},
		{ID: 22, TeamID: page.TeamID, Name: "DB", Type: models.MonitorTypeHTTP, Status: models.MonitorStatusUp},
	}

	incident := models.IncidentWithMonitorIDs{
		Incident: models.Incident{
			ID:        111,
			Status:    models.IncidentStatusDetected,
			Severity:  models.IncidentSeverityMajor,
			IsPublic:  true,
			StartedAt: now,
			CreatedAt: now,
			UpdatedAt: now,
		},
		MonitorIDs: []int64{11, 22},
	}

	mockRepo := &repository.MockRepository{}
	mockRepo.On("StartTransaction", mock.Anything).Return(nil, nil)
	mockRepo.On("DeferRollback", mock.Anything, mock.Anything)
	mockRepo.On("GetStatusPageBySlug", mock.Anything, mock.Anything, "demo").Return(page, nil)
	mockRepo.On("ListStatusPageGroupsByStatusPageID", mock.Anything, mock.Anything, page.ID).Return([]models.StatusPageGroup{}, nil)
	mockRepo.On("ListStatusPageMonitorsByStatusPageID", mock.Anything, mock.Anything, page.ID).Return(monitors, nil)
	mockRepo.On("ListMonitorsByIDs", mock.Anything, mock.Anything, page.TeamID, []int64{11, 22}).Return(monitorRows, nil)
	mockRepo.On("ListPublicIncidentsByMonitorIDs", mock.Anything, mock.Anything, []int64{11, 22}).Return([]models.IncidentWithMonitorIDs{incident}, nil)
	mockRepo.On("ListPublicEventTimelinesByIncidentIDs", mock.Anything, mock.Anything, []int64{int64(111)}).Return([]models.EventTimeline{}, nil)
	mockRepo.On("ListMonitorDailySummaryByMonitorIDs", mock.Anything, mock.Anything, []int64{11, 22}, mock.Anything, mock.Anything).Return([]models.MonitorDailySummary{}, nil)
	mockRepo.On("CommitTransaction", mock.Anything, mock.Anything).Return(nil)

	h := &Handler{Repo: mockRepo}
	c, rec := testutil.NewEchoContext(http.MethodGet, "/status-pages/demo", nil)
	c.SetPath("/status-pages/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("demo")

	err := h.GetPublicStatusPage(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, "Status page returned", resp["message"])

	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)

	incidents, ok := data["incidents"].([]any)
	require.True(t, ok)
	require.Len(t, incidents, 1)

	incidentResp, ok := incidents[0].(map[string]any)
	require.True(t, ok)

	monitorIDs, ok := incidentResp["monitor_id"].([]any)
	require.True(t, ok)
	require.ElementsMatch(t, []any{"11", "22"}, monitorIDs)
}
