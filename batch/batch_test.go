package batch

import (
	"errors"
	"testing"
	"time"

	"github.com/pokt-foundation/transaction-db/types"
	"github.com/sirupsen/logrus"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBatch_RelayBatcher(t *testing.T) {
	c := require.New(t)

	validRelay := types.Relay{
		PoktChainID:              "21",
		EndpointID:               "21",
		SessionKey:               "21",
		PoktNodeAddress:          "21",
		RelayStartDatetime:       time.Now(),
		RelayReturnDatetime:      time.Now(),
		RelayRoundtripTime:       1,
		RelayChainMethodIDs:      []string{"eth_getLogs"},
		RelayDataSize:            21,
		RelayPortalTripTime:      21,
		RelayNodeTripTime:        21,
		RelayURLIsPublicEndpoint: false,
		PortalRegionName:         "region",
		IsAltruistRelay:          false,
		RelaySourceURL:           "example.com",
		PoktNodeDomain:           "node.com",
		PoktNodePublicKey:        "1234",
		RequestID:                "1234",
		PoktTxID:                 "1234",
	}

	invalidRelay := types.Relay{
		EndpointID: "1",
	}

	tests := []struct {
		name        string
		maxSize     int
		maxDuration time.Duration
		timeToWait  time.Duration
		relaysToAdd int
		relayToAdd  types.Relay
		expectedErr error
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
		{
			name:        "Invalid Relay",
			maxSize:     1,
			maxDuration: time.Hour,
			relaysToAdd: 1,
			relayToAdd:  invalidRelay,
			expectedErr: errors.New("PoktChainID is not set"),
		},
	}
	for _, tt := range tests {
		writerMock := &MockRelayWriter{}
		batch := New(tt.maxSize, tt.maxDuration, time.Hour, writerMock, logrus.New())

		writerMock.On("WriteRelays", mock.Anything, mock.Anything).Return(nil).Once()

		for i := 0; i < tt.relaysToAdd; i++ {
			err := batch.AddRelay(tt.relayToAdd)
			c.Equal(tt.expectedErr, err)
		}

		time.Sleep(time.Second)

		c.Equal(0, batch.RelaysSize())
	}
}
