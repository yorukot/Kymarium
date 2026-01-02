package statuspage

import (
	"sort"
	"strconv"

	"github.com/yorukot/kymarium/models"
)

func buildStatusPageElementResponses(groups []models.StatusPageGroup, monitors []models.StatusPageMonitor) []statusPageElementResponse {
	groupMonitors := make(map[int64][]models.StatusPageMonitor, len(groups))
	ungroupedMonitors := make([]models.StatusPageMonitor, 0)
	for _, monitor := range monitors {
		if monitor.GroupID == nil {
			ungroupedMonitors = append(ungroupedMonitors, monitor)
			continue
		}
		groupMonitors[*monitor.GroupID] = append(groupMonitors[*monitor.GroupID], monitor)
	}

	elements := make([]statusPageElementResponse, 0, len(groups)+len(ungroupedMonitors))
	for _, group := range groups {
		monitorList := groupMonitors[group.ID]
		if monitorList == nil {
			monitorList = []models.StatusPageMonitor{}
		}

		element := statusPageElementResponse{
			ID:           strconv.FormatInt(group.ID, 10),
			StatusPageID: strconv.FormatInt(group.StatusPageID, 10),
			Name:         group.Name,
			Type:         group.Type,
			SortOrder:    group.SortOrder,
			Monitor:      false,
			Monitors:     monitorList,
		}

		elements = append(elements, element)
	}

	for _, monitor := range ungroupedMonitors {
		monitorID := strconv.FormatInt(monitor.MonitorID, 10)
		element := statusPageElementResponse{
			ID:           strconv.FormatInt(monitor.ID, 10),
			StatusPageID: strconv.FormatInt(monitor.StatusPageID, 10),
			Name:         monitor.Name,
			Type:         monitor.Type,
			SortOrder:    monitor.SortOrder,
			Monitor:      true,
			MonitorID:    &monitorID,
			Monitors:     []models.StatusPageMonitor{},
		}

		elements = append(elements, element)
	}

	sort.SliceStable(elements, func(i, j int) bool {
		if elements[i].SortOrder == elements[j].SortOrder {
			return elements[i].Name < elements[j].Name
		}
		return elements[i].SortOrder < elements[j].SortOrder
	})

	return elements
}
