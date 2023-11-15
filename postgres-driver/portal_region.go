package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (d *PostgresDriver) WriteRegion(ctx context.Context, region types.PortalRegion) error {
	return d.InsertPortalRegion(ctx, region.PortalRegionName)
}
