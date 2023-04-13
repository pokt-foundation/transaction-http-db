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

func TestBatch_ServiceRecordBatcher(t *testing.T) {
	c := require.New(t)

	validServiceRecord := types.ServiceRecord{
		SessionKey:             "21",
		NodePublicKey:          "21",
		PoktChainID:            "21",
		RequestID:              "21",
		PortalRegionName:       "La Colombia",
		Latency:                21.07,
		Tickets:                2,
		Result:                 "a",
		Available:              true,
		Successes:              21,
		Failures:               7,
		P90SuccessLatency:      21.07,
		MedianSuccessLatency:   21.07,
		WeightedSuccessLatency: 21.07,
		SuccessRate:            21,
	}

	invalidServiceRecord := types.ServiceRecord{
		SessionKey: "1",
	}

	tests := []struct {
		name                string
		maxSize             int
		maxDuration         time.Duration
		timeToWait          time.Duration
		serviceRecordsToAdd int
		serviceRecordToAdd  types.ServiceRecord
		expectedErr         error
	}{
		{
			name:                "Save Relays For Size",
			maxSize:             1,
			maxDuration:         time.Hour,
			serviceRecordsToAdd: 1,
			serviceRecordToAdd:  validServiceRecord,
		},
		{
			name:                "Save Relays For Max Duration Reached",
			maxSize:             2,
			maxDuration:         time.Millisecond,
			serviceRecordsToAdd: 1,
			serviceRecordToAdd:  validServiceRecord,
		},
		{
			name:                "Invalid Relay",
			maxSize:             1,
			maxDuration:         time.Hour,
			serviceRecordsToAdd: 1,
			serviceRecordToAdd:  invalidServiceRecord,
			expectedErr:         errors.New("NodePublicKey is not set"),
		},
	}
	for _, tt := range tests {
		writerMock := &MockServiceRecordWriter{}
		batch := NewServiceRecordBatch(tt.maxSize, tt.maxDuration, time.Hour, writerMock, logrus.New())

		writerMock.On("WriteServiceRecords", mock.Anything, mock.Anything).Return(nil).Once()

		for i := 0; i < tt.serviceRecordsToAdd; i++ {
			err := batch.AddServicRecord(tt.serviceRecordToAdd)
			c.Equal(tt.expectedErr, err)
		}

		time.Sleep(time.Second)

		c.Equal(0, batch.ServiceRecordsSize())
	}
}
