package postgresdriver

import (
	"context"
	"time"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteRelay() {
	tests := []struct {
		name  string
		relay types.Relay
		err   error
	}{
		{
			name: "Success",
			relay: types.Relay{
				PoktChainID:              "21",
				EndpointID:               "21",
				SessionKey:               ts.firstRelay.SessionKey,
				ProtocolAppPublicKey:     "21",
				RelaySourceURL:           "pablo.com",
				PoktNodeAddress:          "21",
				PoktNodeDomain:           "pablos.com",
				PoktNodePublicKey:        "aaa",
				RelayStartDatetime:       time.Now(),
				RelayReturnDatetime:      time.Now(),
				RelayRoundtripTime:       1,
				RelayChainMethodIDs:      []string{"get_height"},
				RelayDataSize:            21,
				RelayPortalTripTime:      21,
				RelayNodeTripTime:        21,
				RelayURLIsPublicEndpoint: false,
				PortalRegionName:         ts.firstRelay.PortalRegionName,
				IsAltruistRelay:          false,
				IsUserRelay:              false,
				RequestID:                "21",
			},
			err: nil,
		},
		{
			name: "Success error relay",
			relay: types.Relay{
				IsError:             true,
				ErrorCode:           21,
				ErrorName:           "favorite number",
				ErrorMessage:        "just Pablo can use it",
				ErrorType:           "chain_check",
				ErrorSource:         "internal",
				PortalRegionName:    ts.firstRelay.PortalRegionName,
				RelayStartDatetime:  time.Now(),
				RelayReturnDatetime: time.Now(),
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			ts.Equal(ts.driver.WriteRelay(context.Background(), tt.relay), tt.err)
		})
	}
}

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteRelays() {
	var relays []*types.Relay
	for i := 0; i < 1000; i++ {
		relays = append(relays, &types.Relay{
			PoktChainID:              "21",
			EndpointID:               "21",
			SessionKey:               ts.firstRelay.SessionKey,
			ProtocolAppPublicKey:     "21",
			RelaySourceURL:           "pablo.com",
			PoktNodeAddress:          "21",
			PoktNodeDomain:           "pablos.com",
			PoktNodePublicKey:        "aaa",
			RelayStartDatetime:       time.Now(),
			RelayReturnDatetime:      time.Now(),
			IsError:                  true,
			ErrorCode:                21,
			ErrorName:                "favorite number",
			ErrorMessage:             "just Pablo can use it",
			ErrorType:                "chain_check",
			ErrorSource:              "internal",
			RelayRoundtripTime:       1,
			RelayChainMethodIDs:      []string{"get_height", "get_balance"},
			RelayDataSize:            21,
			RelayPortalTripTime:      21,
			RelayNodeTripTime:        21,
			RelayURLIsPublicEndpoint: false,
			PortalRegionName:         ts.firstRelay.PortalRegionName,
			IsAltruistRelay:          false,
			IsUserRelay:              false,
			RequestID:                "21",
		})
	}

	tests := []struct {
		name   string
		relays []*types.Relay
		err    error
	}{
		{
			name:   "Success",
			relays: relays,
			err:    nil,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			ts.Equal(ts.driver.WriteRelays(context.Background(), tt.relays), tt.err)
		})
	}
}

func (ts *PGDriverTestSuite) TestPostgresDriver_ReadRelay() {
	tests := []struct {
		name     string
		relayID  int
		expRelay types.Relay
		err      error
	}{
		{
			name:     "Success",
			relayID:  ts.firstRelay.RelayID,
			expRelay: ts.firstRelay,
			err:      nil,
		},
	}
	for _, tt := range tests {
		ts.Run(tt.name, func() {
			relay, err := ts.driver.ReadRelay(context.Background(), tt.relayID)
			ts.Equal(err, tt.err)
			ts.Equal(relay, tt.expRelay)
		})
	}
}
