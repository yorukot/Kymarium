package monitor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/yorukot/knocker/models"
)

type monitorResponse struct {
	ID                string             `json:"id"`
	TeamID            string             `json:"team_id"`
	Name              string             `json:"name"`
	Type              models.MonitorType `json:"type"`
	Config            json.RawMessage    `json:"config"`
	Interval          int                `json:"interval"`
	LastChecked       time.Time          `json:"last_checked"`
	NextCheck         time.Time          `json:"next_check"`
	FailureThreshold  int16              `json:"failure_threshold"`
	RecoveryThreshold int16              `json:"recovery_threshold"`
	NotificationIDs   []string           `json:"notification"`
	UpdatedAt         time.Time          `json:"updated_at"`
	CreatedAt         time.Time          `json:"created_at"`
}

type notificationIDList []int64

func (n *notificationIDList) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("notification must be an array of IDs")
	}

	ids := make([]int64, 0, len(raw))
	for i, item := range raw {
		var str string
		if err := json.Unmarshal(item, &str); err == nil {
			id, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return fmt.Errorf("notification[%d] must be a valid integer string", i)
			}
			if id <= 0 {
				return fmt.Errorf("notification[%d] must be a positive integer", i)
			}
			ids = append(ids, id)
			continue
		}

		var num json.Number
		if err := json.Unmarshal(item, &num); err == nil {
			id, err := num.Int64()
			if err != nil {
				return fmt.Errorf("notification[%d] must be a valid integer", i)
			}
			if id <= 0 {
				return fmt.Errorf("notification[%d] must be a positive integer", i)
			}
			ids = append(ids, id)
			continue
		}

		return fmt.Errorf("notification[%d] must be a string or number", i)
	}

	*n = ids
	return nil
}

func (n notificationIDList) Int64s() []int64 {
	return []int64(n)
}

func newMonitorResponse(m models.Monitor) monitorResponse {
	return monitorResponse{
		ID:                strconv.FormatInt(m.ID, 10),
		TeamID:            strconv.FormatInt(m.TeamID, 10),
		Name:              m.Name,
		Type:              m.Type,
		Config:            m.Config,
		Interval:          m.Interval,
		LastChecked:       m.LastChecked,
		NextCheck:         m.NextCheck,
		FailureThreshold:  m.FailureThreshold,
		RecoveryThreshold: m.RecoveryThreshold,
		NotificationIDs:   formatNotificationIDs(m.NotificationIDs),
		UpdatedAt:         m.UpdatedAt,
		CreatedAt:         m.CreatedAt,
	}
}

func newMonitorResponses(monitors []models.Monitor) []monitorResponse {
	responses := make([]monitorResponse, len(monitors))
	for i, monitor := range monitors {
		responses[i] = newMonitorResponse(monitor)
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
