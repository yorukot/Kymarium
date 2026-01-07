package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/yorukot/kymarium/models"
)

// ListAllRegions fetches all regions from the regions table ordered by ID.
func (r *PGRepository) ListAllRegions(ctx context.Context, tx pgx.Tx) ([]models.Region, error) {
	const query = `
		SELECT id, name
		FROM regions
		ORDER BY id
	`

	var regions []models.Region
	if err := pgxscan.Select(ctx, tx, &regions, query); err != nil {
		return nil, err
	}

	return regions, nil
}
