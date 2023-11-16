package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteRegion() {
	tests := []struct {
		name   string
		region types.PortalRegion
		err    error
	}{
		{
			name: "Success",
			region: types.PortalRegion{
				PortalRegionName: "Cartago",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		ts.Equal(ts.driver.WriteRegion(context.Background(), tt.region), tt.err)
	}
}
