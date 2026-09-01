package ingest

import (
	"log/slog"
	"sync/atomic"

	"github.com/rafawastaken/tick-storm/backend/internal/client/binance"
	"github.com/rafawastaken/tick-storm/backend/internal/crypto"
)

type Worker struct {
	client  *binance.Client
	svc     *crypto.Service
	log     *slog.Logger
	symbols []string

	ingested atomic.Int64
	skipped  atomic.Int64
	// Stored in microseconds and divided on display: whole milliseconds
	// round a sub-millisecond lag down to zero.
	lagUS atomic.Int64
}

func NewWorker(client *binance.Client, svc *crypto.Service, log *slog.Logger, symbols []string) *Worker {
	return &Worker{
		client:  client,
		svc:     svc,
		log:     log,
		symbols: symbols,
	}
}
