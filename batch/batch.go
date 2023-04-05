package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pokt-foundation/transaction-db/types"
	"github.com/sirupsen/logrus"
)

type RelayWriter interface {
	WriteRelays(ctx context.Context, relays []types.Relay) error
}

type Batch struct {
	relays      []types.Relay
	rwMutex     sync.RWMutex
	batchChan   chan types.Relay
	maxSize     int
	maxDuration time.Duration
	timeoutDB   time.Duration
	writer      RelayWriter
	log         *logrus.Logger
}

func (b *Batch) logError(err error) {
	fields := logrus.Fields{
		"err": err.Error(),
	}

	b.log.WithFields(fields).Error(err)
}

func New(maxSize int, maxDuration, timeoutDB time.Duration, writer RelayWriter, logger *logrus.Logger) *Batch {
	batch := &Batch{
		maxSize:     maxSize,
		maxDuration: maxDuration,
		timeoutDB:   timeoutDB,
		writer:      writer,
		batchChan:   make(chan types.Relay, 32),
		log:         logger,
	}

	go batch.RelayBatcher()

	return batch
}

func (b *Batch) AddRelay(relay types.Relay) error {
	if err := relay.Validate(); err != nil {
		return err
	}

	b.batchChan <- relay

	return nil
}

func (b *Batch) RelaysSize() int {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	return len(b.relays)
}

func (b *Batch) addRelay(relay types.Relay) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.relays = append(b.relays, relay)
}

func (b *Batch) RelayBatcher() {
	ticker := time.NewTicker(b.maxDuration)
	defer ticker.Stop()

	for {
		select {
		case relay := <-b.batchChan:
			b.log.Debug("Relay received in batch")
			b.addRelay(relay)

			if b.RelaysSize() >= b.maxSize {
				b.log.Debug("Max size on relay batcher reached")
				if err := b.SaveRelaysToDB(); err != nil {
					b.logError(fmt.Errorf("error saving batch: %s", err))
				}
				// Reset the ticker when max size is reached
				ticker = time.NewTicker(b.maxDuration)
			}

		case <-ticker.C:
			b.log.Debug("Max duration on relay batcher reached")
			if err := b.SaveRelaysToDB(); err != nil {
				b.logError(fmt.Errorf("error saving batch: %s", err))
			}
		}
	}
}

func (b *Batch) SaveRelaysToDB() error {
	b.rwMutex.Lock()
	relays := make([]types.Relay, len(b.relays))
	copy(relays, b.relays)
	b.relays = nil
	b.rwMutex.Unlock()

	if len(relays) == 0 {
		b.log.Info("No relay was saved")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.timeoutDB)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- b.writer.WriteRelays(ctx, relays)
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
