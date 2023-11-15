package postgresdriver

import (
	"context"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteServiceRecord() {
	tests := []struct {
		name          string
		serviceRecord types.ServiceRecord
		err           error
	}{
		{
			name: "Success",
			serviceRecord: types.ServiceRecord{
				NodePublicKey:          "21",
				PoktChainID:            "21",
				SessionKey:             ts.firstServiceRecord.SessionKey,
				RequestID:              "21",
				PortalRegionName:       ts.firstServiceRecord.PortalRegionName,
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
	}
	for _, tt := range tests {
		ts.Equal(ts.driver.WriteServiceRecord(context.Background(), tt.serviceRecord), tt.err)
	}
}

func (ts *PGDriverTestSuite) TestPostgresDriver_WriteServiceRecords() {
	var serviceRecords []*types.ServiceRecord
	for i := 0; i < 1000; i++ {
		serviceRecords = append(serviceRecords, &types.ServiceRecord{
			NodePublicKey:          "21",
			PoktChainID:            "21",
			SessionKey:             ts.firstServiceRecord.SessionKey,
			RequestID:              "21",
			PortalRegionName:       ts.firstServiceRecord.PortalRegionName,
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
		})
	}

	tests := []struct {
		name           string
		serviceRecords []*types.ServiceRecord
		err            error
	}{
		{
			name:           "Success",
			serviceRecords: serviceRecords,
			err:            nil,
		},
	}
	for _, tt := range tests {
		ts.Equal(ts.driver.WriteServiceRecords(context.Background(), tt.serviceRecords), tt.err)
	}
}

func (ts *PGDriverTestSuite) TestPostgresDriver_ReadServiceRecord() {
	tests := []struct {
		name             string
		serviceRecordID  int
		expServiceRecord types.ServiceRecord
		err              error
	}{
		{
			name:             "Success",
			serviceRecordID:  ts.firstServiceRecord.ServiceRecordID,
			expServiceRecord: ts.firstServiceRecord,
			err:              nil,
		},
	}
	for _, tt := range tests {
		serviceRecord, err := ts.driver.ReadServiceRecord(context.Background(), tt.serviceRecordID)
		ts.Equal(err, tt.err)
		ts.Equal(serviceRecord, tt.expServiceRecord)
	}
}
