package handler

import (
	"context"
	"time"

	"github.com/yorukot/kymarium/models"
	"github.com/yorukot/kymarium/repository"
	"go.uber.org/zap"
)

const (
	defaultPingFlushSize     = 1000
	defaultPingFlushInterval = 1 * time.Second
)

// PingRecorder buffers ping results and periodically writes them to the database.
type PingRecorder struct {
	repo          repository.Repository
	buffer        chan models.Ping
	flushSize     int
	flushInterval time.Duration
	done          chan struct{}
	stopped       chan struct{}
}

// NewPingRecorder creates a PingRecorder with default buffering settings.
// NewPingRecorder creates a ping recorder with default buffering settings.
func NewPingRecorder(repo repository.Repository) *PingRecorder {
	recorder := &PingRecorder{
		repo:          repo,
		buffer:        make(chan models.Ping, defaultPingFlushSize*4),
		flushSize:     defaultPingFlushSize,
		flushInterval: defaultPingFlushInterval,
		done:          make(chan struct{}),
		stopped:       make(chan struct{}),
	}

	go recorder.run()
	return recorder
}

// Record enqueues a ping for batch persistence.
func (r *PingRecorder) Record(ctx context.Context, ping models.Ping) {
	select {
	case r.buffer <- ping:
	default:
		// Prevent handler backpressure from stalling task processing.
		go r.flushOne(ping)
	}
}

func (r *PingRecorder) run() {
	ticker := time.NewTicker(r.flushInterval)
	defer ticker.Stop()

	batch := make([]models.Ping, 0, r.flushSize)
	flushThreshold := int(float64(r.flushSize) * 0.8)

	for {
		select {
		case ping := <-r.buffer:
			batch = append(batch, ping)
			if len(batch) >= flushThreshold {
				batch = r.flush(batch)
			}
		case <-ticker.C:
			batch = r.flush(batch)
		}
	}
}

func (r *PingRecorder) flushOne(ping models.Ping) {
	if remaining := r.flush([]models.Ping{ping}); len(remaining) > 0 {
		// If the flush failed, fall back to a blocking enqueue to avoid losing the record.
		r.buffer <- ping
	}
}

func (r *PingRecorder) flush(batch []models.Ping) []models.Ping {
	if len(batch) == 0 {
		return batch
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.repo.StartTransaction(ctx)
	if err != nil {
		zap.L().Error("failed to start transaction for ping batch", zap.Error(err))
		return batch
	}
	defer r.repo.DeferRollback(ctx, tx)

	if err := r.repo.BatchInsertPings(ctx, tx, batch); err != nil {
		zap.L().Error("failed to insert ping batch", zap.Int("count", len(batch)), zap.Error(err))
		return batch
	}

	if err := r.repo.CommitTransaction(ctx, tx); err != nil {
		return batch
	}

	zap.L().Debug("flushed ping batch", zap.Int("count", len(batch)))
	return batch[:0]
}
