package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteSession() {
	tests := []struct {
		name    string
		session types.PocketSession
		err     error
	}{
		{
			name: "Success",
			session: types.PocketSession{
				SessionKey:       "21",
				SessionHeight:    21,
				PortalRegionName: "europe-southwest1",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		ts.Equal(ts.driver.WriteSession(context.Background(), tt.session), tt.err)
		ts.Equal(ts.driver.WriteSession(context.Background(), tt.session), types.ErrRepeatedSessionKey)
	}
}
