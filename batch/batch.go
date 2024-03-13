package batch

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pokt-foundation/portal-middleware/metrics/exporter"
	"github.com/pokt-foundation/transaction-http-db/metric"
	"go.uber.org/zap"
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
	log         *zap.Logger
	index       atomic.Int32
	errCounter  exporter.Counter
	dataCounter exporter.Counter
}

func (b *Batch[T]) logError(err error) {
	b.log.Error(err.Error(), zap.String("err", err.Error()), zap.String("name", b.name))
}

func NewBatch[T Validator](maxSize, chanSize int, name string, maxDuration, timeoutDB time.Duration, writer writerFunc[T], logger *zap.Logger, metrics exporter.MetricExporter) *Batch[T] {

	errCounter := metrics.Counter(metric.CategorySubscription, metric.NameError)
	dataCounter := metrics.Counter(metric.CategorySubscription, metric.NameData)
	batch := &Batch[T]{
		maxSize:     maxSize,
		name:        name,
		maxDuration: maxDuration,
		timeoutDB:   timeoutDB,
		writer:      writer,
		batchChan:   make(chan T, chanSize),
		log:         logger,
		items:       make([]T, maxSize),
		index:       atomic.Int32{},
		dataCounter: dataCounter,
		errCounter:  errCounter,
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
	return int(b.index.Load())
}

func (b *Batch[T]) add(item T) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.items[b.index.Load()] = item
	b.index.Add(1)
}

func (b *Batch[T]) Batcher() {
	ticker := time.NewTicker(b.maxDuration)
	defer ticker.Stop()

	for {
		select {
		case item := <-b.batchChan:
			b.log.Debug(fmt.Sprintf("item received in %s batch", b.name))
			b.add(item)

			if b.Size() >= b.maxSize {
				b.log.Debug(fmt.Sprintf("max size on %s batcher reached", b.name))
				if err := b.Save(); err != nil {
					b.logError(fmt.Errorf("error saving %s batch: %s", b.name, err))
					b.reportErrorMetric()
				}
				// Reset the ticker when max size is reached
				ticker = time.NewTicker(b.maxDuration)
			}

		case <-ticker.C:
			b.log.Debug(fmt.Sprintf("max duration on %s batcher reached", b.name))
			if err := b.Save(); err != nil {
				b.logError(fmt.Errorf("error saving %s batch: %s", b.name, err))
				b.reportErrorMetric()
			}
		}
	}
}

func (b *Batch[T]) reportErrorMetric() {
	b.dataCounter.Add(fmt.Sprintf("%s_lost", b.name), float64(b.Size()))
	b.errCounter.Inc(fmt.Sprintf("%s_batch_failed", b.name))
}

func (b *Batch[T]) Save() error {
	b.rwMutex.Lock()

	size := b.index.Load()
	items := make([]T, size)
	copy(items, b.items[:size])
	b.index.Store(0)

	b.rwMutex.Unlock()

	if len(items) == 0 {
		b.log.Warn(fmt.Sprintf("no item was saved on %s", b.name))
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

	b.dataCounter.Add(fmt.Sprintf("%s_saved", b.name), float64(size)) // the name will change from relay and count
	return nil
}
