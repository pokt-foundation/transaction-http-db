package batch

import (
	context "context"
	"fmt"
	"sync"
	"time"

	types "github.com/pokt-foundation/transaction-db/types"
	"github.com/sirupsen/logrus"
)

type ServiceRecordWriter interface {
	WriteServiceRecords(ctx context.Context, serviceRecords []types.ServiceRecord) error
}

type ServiceRecordBatch struct {
	serviceRecords []types.ServiceRecord
	rwMutex        sync.RWMutex
	batchChan      chan types.ServiceRecord
	maxSize        int
	maxDuration    time.Duration
	timeoutDB      time.Duration
	writer         ServiceRecordWriter
	log            *logrus.Logger
}

func (b *ServiceRecordBatch) logError(err error) {
	fields := logrus.Fields{
		"err": err.Error(),
	}

	b.log.WithFields(fields).Error(err)
}

func NewServiceRecordBatch(maxSize int, maxDuration, timeoutDB time.Duration, writer ServiceRecordWriter, logger *logrus.Logger) *ServiceRecordBatch {
	batch := &ServiceRecordBatch{
		maxSize:     maxSize,
		maxDuration: maxDuration,
		timeoutDB:   timeoutDB,
		writer:      writer,
		batchChan:   make(chan types.ServiceRecord, 32),
		log:         logger,
	}

	go batch.ServiceRecordBatcher()

	return batch
}

func (b *ServiceRecordBatch) AddServicRecord(serviceRecord types.ServiceRecord) error {
	if err := serviceRecord.Validate(); err != nil {
		return err
	}

	b.batchChan <- serviceRecord

	return nil
}

func (b *ServiceRecordBatch) ServiceRecordsSize() int {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	return len(b.serviceRecords)
}

func (b *ServiceRecordBatch) addServiceRecord(serviceRecord types.ServiceRecord) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.serviceRecords = append(b.serviceRecords, serviceRecord)
}

func (b *ServiceRecordBatch) ServiceRecordBatcher() {
	ticker := time.NewTicker(b.maxDuration)
	defer ticker.Stop()

	for {
		select {
		case serviceRecord := <-b.batchChan:
			b.log.Debug("Service record received in batch")
			b.addServiceRecord(serviceRecord)

			if b.ServiceRecordsSize() >= b.maxSize {
				b.log.Debug("Max size on service record batcher reached")
				if err := b.SaveServiceRecordsToDB(); err != nil {
					b.logError(fmt.Errorf("error saving batch: %s", err))
				}
				// Reset the ticker when max size is reached
				ticker = time.NewTicker(b.maxDuration)
			}

		case <-ticker.C:
			b.log.Debug("Max duration on service record batcher reached")
			if err := b.SaveServiceRecordsToDB(); err != nil {
				b.logError(fmt.Errorf("error saving batch: %s", err))
			}
		}
	}
}

func (b *ServiceRecordBatch) SaveServiceRecordsToDB() error {
	b.rwMutex.Lock()
	serviceRecords := make([]types.ServiceRecord, len(b.serviceRecords))
	copy(serviceRecords, b.serviceRecords)
	b.serviceRecords = nil
	b.rwMutex.Unlock()

	if len(serviceRecords) == 0 {
		b.log.Info("No service record was saved")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.timeoutDB)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- b.writer.WriteServiceRecords(ctx, serviceRecords)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
