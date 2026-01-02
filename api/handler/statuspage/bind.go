package statuspage

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yorukot/kymarium/models"
)

func bindStatusPageUpsert(c echo.Context) (statusPageUpsertRequest, error) {
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if isFormContentType(contentType) {
		if err := c.Request().ParseForm(); err != nil {
			return statusPageUpsertRequest{}, err
		}
		return parseStatusPageUpsertForm(c.Request().Form)
	}

	var req statusPageUpsertRequest
	if err := c.Bind(&req); err != nil {
		return statusPageUpsertRequest{}, err
	}
	return req, nil
}

func isFormContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data")
}

func parseStatusPageUpsertForm(values url.Values) (statusPageUpsertRequest, error) {
	req := statusPageUpsertRequest{
		Title: firstNonEmpty(values, "name", "title"),
		Slug:  firstNonEmpty(values, "slug"),
	}

	elementInputs := make(map[int]*statusPageElementInput)
	elementMonitors := make(map[int]map[int]*statusPageMonitorInput)
	groupInputs := make(map[int]*statusPageGroupInput)
	monitorInputs := make(map[int]*statusPageMonitorInput)

	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		value := vals[len(vals)-1]
		if value == "" {
			continue
		}

		switch {
		case strings.HasPrefix(key, "elements."):
			parseElementFormField(elementInputs, elementMonitors, key, value)
		case strings.HasPrefix(key, "groups."):
			parseGroupFormField(groupInputs, key, value)
		case strings.HasPrefix(key, "monitors."):
			parseMonitorFormField(monitorInputs, key, value)
		}
	}

	req.Elements = buildElementInputs(elementInputs, elementMonitors)
	req.Groups = buildGroupInputs(groupInputs)
	req.Monitors = buildMonitorInputs(monitorInputs)

	return req, nil
}

func parseElementFormField(
	elements map[int]*statusPageElementInput,
	elementMonitors map[int]map[int]*statusPageMonitorInput,
	key string,
	value string,
) {
	parts := strings.Split(key, ".")
	if len(parts) < 3 {
		return
	}

	elementIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	element := ensureElementInput(elements, elementIndex)
	field := parts[2]

	if field == "monitors" {
		if len(parts) < 5 {
			return
		}
		monitorIndex, err := strconv.Atoi(parts[3])
		if err != nil {
			return
		}
		monitor := ensureElementMonitorInput(elementMonitors, elementIndex, monitorIndex)
		setMonitorField(monitor, parts[4], value)
		return
	}

	setElementField(element, field, value)
}

func parseGroupFormField(groups map[int]*statusPageGroupInput, key string, value string) {
	parts := strings.Split(key, ".")
	if len(parts) < 3 {
		return
	}

	groupIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	group := ensureGroupInput(groups, groupIndex)
	setGroupField(group, parts[2], value)
}

func parseMonitorFormField(monitors map[int]*statusPageMonitorInput, key string, value string) {
	parts := strings.Split(key, ".")
	if len(parts) < 3 {
		return
	}

	monitorIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	monitor := ensureMonitorInput(monitors, monitorIndex)
	setMonitorField(monitor, parts[2], value)
}

func ensureElementInput(elements map[int]*statusPageElementInput, index int) *statusPageElementInput {
	element, ok := elements[index]
	if !ok {
		element = &statusPageElementInput{}
		elements[index] = element
	}
	return element
}

func ensureElementMonitorInput(
	elementMonitors map[int]map[int]*statusPageMonitorInput,
	elementIndex int,
	monitorIndex int,
) *statusPageMonitorInput {
	monitorMap, ok := elementMonitors[elementIndex]
	if !ok {
		monitorMap = make(map[int]*statusPageMonitorInput)
		elementMonitors[elementIndex] = monitorMap
	}
	monitor, ok := monitorMap[monitorIndex]
	if !ok {
		monitor = &statusPageMonitorInput{}
		monitorMap[monitorIndex] = monitor
	}
	return monitor
}

func ensureGroupInput(groups map[int]*statusPageGroupInput, index int) *statusPageGroupInput {
	group, ok := groups[index]
	if !ok {
		group = &statusPageGroupInput{}
		groups[index] = group
	}
	return group
}

func ensureMonitorInput(monitors map[int]*statusPageMonitorInput, index int) *statusPageMonitorInput {
	monitor, ok := monitors[index]
	if !ok {
		monitor = &statusPageMonitorInput{}
		monitors[index] = monitor
	}
	return monitor
}

func setElementField(element *statusPageElementInput, field string, value string) {
	switch field {
	case "id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			element.ID = &parsed
		}
	case "name":
		element.Name = value
	case "type":
		element.Type = models.StatusPageElementType(value)
	case "sortOrder", "sort_order":
		if parsed, err := strconv.Atoi(value); err == nil {
			element.SortOrder = parsed
		}
	case "monitor":
		element.Monitor = parseBool(value)
	case "monitorId", "monitor_id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			element.MonitorID = &parsed
		}
	}
}

func setGroupField(group *statusPageGroupInput, field string, value string) {
	switch field {
	case "id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			group.ID = &parsed
		}
	case "name":
		group.Name = value
	case "type":
		group.Type = models.StatusPageElementType(value)
	case "sortOrder", "sort_order":
		if parsed, err := strconv.Atoi(value); err == nil {
			group.SortOrder = parsed
		}
	}
}

func setMonitorField(monitor *statusPageMonitorInput, field string, value string) {
	switch field {
	case "id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			monitor.ID = &parsed
		}
	case "monitorId", "monitor_id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			monitor.MonitorID = parsed
		}
	case "groupId", "group_id":
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			monitor.GroupID = &parsed
		}
	case "name":
		monitor.Name = value
	case "type":
		monitor.Type = models.StatusPageElementType(value)
	case "sortOrder", "sort_order":
		if parsed, err := strconv.Atoi(value); err == nil {
			monitor.SortOrder = parsed
		}
	}
}

func parseBool(value string) bool {
	switch strings.ToLower(value) {
	case "true", "1", "on", "yes":
		return true
	default:
		return false
	}
}

func buildElementInputs(
	elements map[int]*statusPageElementInput,
	elementMonitors map[int]map[int]*statusPageMonitorInput,
) []statusPageElementInput {
	indices := make([]int, 0, len(elements))
	for index := range elements {
		indices = append(indices, index)
	}
	sortInts(indices)

	result := make([]statusPageElementInput, 0, len(indices))
	for _, index := range indices {
		element := *elements[index]
		if monitorMap, ok := elementMonitors[index]; ok {
			element.Monitors = buildMonitorInputs(monitorMap)
		}
		result = append(result, element)
	}
	return result
}

func buildGroupInputs(groups map[int]*statusPageGroupInput) []statusPageGroupInput {
	indices := make([]int, 0, len(groups))
	for index := range groups {
		indices = append(indices, index)
	}
	sortInts(indices)

	result := make([]statusPageGroupInput, 0, len(indices))
	for _, index := range indices {
		result = append(result, *groups[index])
	}
	return result
}

func buildMonitorInputs(monitors map[int]*statusPageMonitorInput) []statusPageMonitorInput {
	indices := make([]int, 0, len(monitors))
	for index := range monitors {
		indices = append(indices, index)
	}
	sortInts(indices)

	result := make([]statusPageMonitorInput, 0, len(indices))
	for _, index := range indices {
		result = append(result, *monitors[index])
	}
	return result
}

func sortInts(values []int) {
	if len(values) < 2 {
		return
	}
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}

func firstNonEmpty(values url.Values, keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(values.Get(key))
		if value != "" {
			return value
		}
	}
	return ""
}
