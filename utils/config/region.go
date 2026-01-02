package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yorukot/kymarium/models"
)

// InitRegionConfig loads regions from the database once and caches them.
// name matches legacy AppRegion/AppRegions, IDs match AppRegionID/AppRegionIDs.
func InitRegionConfig(pool *pgxpool.Pool) ([]models.Region, error) {
	ctx := context.Background()
	regionsOnce.Do(func() {
		rows, err := pool.Query(ctx, `SELECT id, name, display_name FROM regions ORDER BY id`)
		if err != nil {
			regionsErr = fmt.Errorf("list regions: %w", err)
			return
		}
		defer rows.Close()

		byID := make(map[int64]models.Region)
		byKey := make(map[string]models.Region)

		for rows.Next() {
			var r models.Region
			if err := rows.Scan(&r.ID, &r.Name, &r.DisplayName); err != nil {
				regionsErr = fmt.Errorf("scan region: %w", err)
				return
			}
			byID[r.ID] = r
			if key := r.Name; key != "" {
				byKey[key] = r
			}
		}

		if err := rows.Err(); err != nil {
			regionsErr = fmt.Errorf("iterate regions: %w", err)
			return
		}

		regionsByID = byID
		regionsByKey = byKey
	})

	return Regions(), regionsErr
}

// Regions returns cached regions sorted by ID; empty if not initialized.
func Regions() []models.Region {
	result := make([]models.Region, 0, len(regionsByID))
	for _, r := range regionsByID {
		result = append(result, r)
	}
	return result
}

// RegionByID looks up a cached region by ID.
func RegionByID(id int64) models.Region {
	r, _ := regionsByID[id]
	return r
}

// RegionByName looks up a cached region by its name (case-insensitive).
func RegionByName(name string) models.Region {
	r, _ := regionsByKey[name]
	return r
}
