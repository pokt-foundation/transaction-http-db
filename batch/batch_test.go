package batch

import (
	"testing"
	"time"

	"github.com/pokt-foundation/transaction-db/types"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestBatch_RelayBatcher(t *testing.T) {
	c := require.New(t)

	validRelay := types.Relay{
		PoktChainID:              "21",
		EndpointID:               "21",
		SessionKey:               "21",
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
		RelayChainMethodIDs:      []string{"get_height"},
		RelayDataSize:            21,
		RelayPortalTripTime:      21,
		RelayNodeTripTime:        21,
		RelayURLIsPublicEndpoint: false,
		PortalRegionName:         "La Colombia",
		IsAltruistRelay:          false,
		IsUserRelay:              false,
		RequestID:                "21",
		PoktTxID:                 "21",
	}

	tests := []struct {
		name        string
		maxSize     int
		maxDuration time.Duration
		timeToWait  time.Duration
		relaysToAdd int
		relayToAdd  types.Relay
	}{
		{
			name:        "Save Relays For Size",
			maxSize:     1,
			maxDuration: time.Hour,
			relaysToAdd: 1,
			relayToAdd:  validRelay,
		},
		{
			name:        "Save Relays For Max Duration Reached",
			maxSize:     2,
			maxDuration: time.Millisecond,
			relaysToAdd: 1,
			relayToAdd:  validRelay,
		},
	}
	for _, tt := range tests {
		writerMock := &MockRelayWriter{}
		batch := NewBatch(tt.maxSize, 21, "relay", tt.maxDuration, time.Hour, writerMock.WriteRelays, zap.NewNop())

		writerMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()

		for i := 0; i < tt.relaysToAdd; i++ {
			err := batch.Add(&tt.relayToAdd)
			c.NoError(err)
		}

		time.Sleep(time.Second)

		c.Equal(0, batch.Size())
	}
}
