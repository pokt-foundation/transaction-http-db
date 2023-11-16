package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceRecord_ValidateStruct(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name          string
		serviceRecord ServiceRecord
		err           error
	}{
		{
			name: "Success service record",
			serviceRecord: ServiceRecord{
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
			},
			err: nil,
		},
		{
			name: "Success service record without optional fields",
			serviceRecord: ServiceRecord{
				SessionKey:             "21",
				NodePublicKey:          "21",
				PoktChainID:            "21",
				RequestID:              "21",
				PortalRegionName:       "La Colombia",
				Tickets:                2,
				Available:              true,
				P90SuccessLatency:      21.07,
				MedianSuccessLatency:   21.07,
				WeightedSuccessLatency: 21.07,
				SuccessRate:            21,
			},
			err: nil,
		},
		{
			name: "Failure service record id set",
			serviceRecord: ServiceRecord{
				ServiceRecordID:        21,
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
			},
			err: errors.New("ServiceRecordID should not be set"),
		},
		{
			name: "Failure service record field not set",
			serviceRecord: ServiceRecord{
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
			},
			err: errors.New("SessionKey is not set"),
		},
	}

	for _, tt := range tests {
		c.Equal(tt.err, tt.serviceRecord.Validate())
	}
}
