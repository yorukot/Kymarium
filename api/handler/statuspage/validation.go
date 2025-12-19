package statuspage

import (
	"fmt"
)

// validateStatusPagePayload enforces sort order uniqueness and group references.
func validateStatusPagePayload(req statusPageUpsertRequest) error {
	groupOrders := make(map[int]struct{})
	groupIDs := make(map[int64]struct{})

	for _, g := range req.Groups {
		if _, exists := groupOrders[g.SortOrder]; exists {
			return fmt.Errorf("duplicate group sort_order %d", g.SortOrder)
		}
		groupOrders[g.SortOrder] = struct{}{}
		if g.ID != nil {
			groupIDs[*g.ID] = struct{}{}
		}
	}

	monitorOrders := make(map[string]struct{}) // key: groupID|nil:sortOrder
	for _, m := range req.Monitors {
		groupKey := "nil"
		if m.GroupID != nil {
			if _, ok := groupIDs[*m.GroupID]; !ok {
				return fmt.Errorf("monitor group_id %d not present in groups list", *m.GroupID)
			}
			groupKey = fmt.Sprintf("%d", *m.GroupID)
		}

		key := fmt.Sprintf("%s:%d", groupKey, m.SortOrder)
		if _, exists := monitorOrders[key]; exists {
			return fmt.Errorf("duplicate monitor sort_order %d within group %s", m.SortOrder, groupKey)
		}
		monitorOrders[key] = struct{}{}
	}

	return nil
}
