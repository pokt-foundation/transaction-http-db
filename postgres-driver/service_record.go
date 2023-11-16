package postgresdriver

import (
	"context"
	"time"

	"github.com/pokt-foundation/transaction-http-db/types"
)

func (d *PostgresDriver) WriteServiceRecord(ctx context.Context, serviceRecord types.ServiceRecord) error {
	createdAt := time.Now()

	_, err := d.InsertServiceRecords(ctx, []InsertServiceRecordsParams{
		{
			NodePublicKey:          serviceRecord.NodePublicKey,
			PoktChainID:            serviceRecord.PoktChainID,
			SessionKey:             serviceRecord.SessionKey,
			RequestID:              serviceRecord.RequestID,
			PortalRegionName:       serviceRecord.PortalRegionName,
			Latency:                serviceRecord.Latency,
			Tickets:                int32(serviceRecord.Tickets),
			Result:                 serviceRecord.Result,
			Available:              serviceRecord.Available,
			Successes:              int32(serviceRecord.Successes),
			Failures:               int32(serviceRecord.Failures),
			P90SuccessLatency:      serviceRecord.P90SuccessLatency,
			MedianSuccessLatency:   serviceRecord.MedianSuccessLatency,
			WeightedSuccessLatency: serviceRecord.WeightedSuccessLatency,
			SuccessRate:            serviceRecord.SuccessRate,
			CreatedAt:              newTimestamp(createdAt),
			UpdatedAt:              newTimestamp(createdAt),
		},
	})

	return err
}

func (d *PostgresDriver) WriteServiceRecords(ctx context.Context, serviceRecords []*types.ServiceRecord) error {
	createdAt := time.Now()

	serviceRecordParams := make([]InsertServiceRecordsParams, 0, len(serviceRecords))

	for _, serviceRecord := range serviceRecords {
		serviceRecordParams = append(serviceRecordParams, InsertServiceRecordsParams{
			NodePublicKey:          serviceRecord.NodePublicKey,
			PoktChainID:            serviceRecord.PoktChainID,
			SessionKey:             serviceRecord.SessionKey,
			RequestID:              serviceRecord.RequestID,
			PortalRegionName:       serviceRecord.PortalRegionName,
			Latency:                serviceRecord.Latency,
			Tickets:                int32(serviceRecord.Tickets),
			Result:                 serviceRecord.Result,
			Available:              serviceRecord.Available,
			Successes:              int32(serviceRecord.Successes),
			Failures:               int32(serviceRecord.Failures),
			P90SuccessLatency:      serviceRecord.P90SuccessLatency,
			MedianSuccessLatency:   serviceRecord.MedianSuccessLatency,
			WeightedSuccessLatency: serviceRecord.WeightedSuccessLatency,
			SuccessRate:            serviceRecord.SuccessRate,
			CreatedAt:              newTimestamp(createdAt),
			UpdatedAt:              newTimestamp(createdAt),
		})
	}

	_, err := d.InsertServiceRecords(ctx, serviceRecordParams)

	return err
}

func (d *PostgresDriver) ReadServiceRecord(ctx context.Context, serviceRecordID int) (types.ServiceRecord, error) {
	serviceRecord, err := d.SelectServiceRecord(ctx, int64(serviceRecordID))
	if err != nil {
		return types.ServiceRecord{}, err
	}

	return types.ServiceRecord{
		NodePublicKey:          serviceRecord.NodePublicKey,
		PoktChainID:            serviceRecord.PoktChainID,
		ServiceRecordID:        int(serviceRecord.ID),
		SessionKey:             serviceRecord.SessionKey,
		RequestID:              serviceRecord.RequestID,
		PortalRegionName:       serviceRecord.PortalRegionName,
		Latency:                serviceRecord.Latency,
		Tickets:                int(serviceRecord.Tickets),
		Result:                 serviceRecord.Result,
		Available:              serviceRecord.Available,
		Successes:              int(serviceRecord.Successes),
		Failures:               int(serviceRecord.Failures),
		P90SuccessLatency:      serviceRecord.P90SuccessLatency,
		MedianSuccessLatency:   serviceRecord.MedianSuccessLatency,
		WeightedSuccessLatency: serviceRecord.WeightedSuccessLatency,
		SuccessRate:            serviceRecord.SuccessRate,
		CreatedAt:              serviceRecord.CreatedAt.Time,
		UpdatedAt:              serviceRecord.UpdatedAt.Time,
	}, nil
}
