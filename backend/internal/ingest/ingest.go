package ingest

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/rafawastaken/tick-storm/backend/internal/client/binance"
	"github.com/rafawastaken/tick-storm/backend/internal/crypto"
)

type Worker struct {
	client  *binance.Client
	svc     *crypto.Service
	log     *slog.Logger
	symbols []string

	ingested    atomic.Int64
	skipped     atomic.Int64
	lagUS       atomic.Int64
	lastTradeAt atomic.Int64
}

func NewWorker(client *binance.Client, svc *crypto.Service, log *slog.Logger, symbols []string) *Worker {
	return &Worker{
		client:  client,
		svc:     svc,
		log:     log,
		symbols: symbols,
	}
}

func (w *Worker) SilentFor() time.Duration {
	return time.Since(time.Unix(0, w.lastTradeAt.Load()))
}
