package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Validator interface {
	Validate() error
}

type writerFunc[T Validator] func(context.Context, []T) error

type Batch[T Validator] struct {
	items       []T
	rwMutex     sync.RWMutex
	batchChan   chan T
	maxSize     int
	name        string
	maxDuration time.Duration
	timeoutDB   time.Duration
	writer      writerFunc[T]
	log         *logrus.Logger
}

func (b *Batch[T]) logError(err error) {
	fields := logrus.Fields{
		"err":  err.Error(),
		"name": b.name,
	}

	b.log.WithFields(fields).Error(err)
}

func NewBatch[T Validator](maxSize int, name string, maxDuration, timeoutDB time.Duration, writer writerFunc[T], logger *logrus.Logger) *Batch[T] {
	batch := &Batch[T]{
		maxSize:     maxSize,
		name:        name,
		maxDuration: maxDuration,
		timeoutDB:   timeoutDB,
		writer:      writer,
		batchChan:   make(chan T, 32),
		log:         logger,
	}

	go batch.Batcher()

	return batch
}

func (b *Batch[T]) Add(item T) error {
	if err := item.Validate(); err != nil {
		return err
	}

	b.batchChan <- item

	return nil
}

func (b *Batch[T]) Size() int {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	return len(b.items)
}

func (b *Batch[T]) add(item T) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.items = append(b.items, item)
}

func (b *Batch[T]) Batcher() {
	ticker := time.NewTicker(b.maxDuration)
	defer ticker.Stop()

	for {
		select {
		case item := <-b.batchChan:
			b.log.Debugf("item received in %s batch", b.name)
			b.add(item)

			if b.Size() >= b.maxSize {
				b.log.Debugf("max size on %s batcher reached", b.name)
				if err := b.Save(); err != nil {
					b.logError(fmt.Errorf("error saving %s batch: %s", b.name, err))
				}
				// Reset the ticker when max size is reached
				ticker = time.NewTicker(b.maxDuration)
			}

		case <-ticker.C:
			b.log.Debugf("max duration on %s batcher reached", b.name)
			if err := b.Save(); err != nil {
				b.logError(fmt.Errorf("error saving %s batch: %s", b.name, err))
			}
		}
	}
}

func (b *Batch[T]) Save() error {
	b.rwMutex.Lock()
	items := make([]T, len(b.items))
	copy(items, b.items)
	b.items = nil
	b.rwMutex.Unlock()

	if len(items) == 0 {
		b.log.Infof("no item was saved on %s", b.name)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.timeoutDB)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- b.writer(ctx, items)
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
