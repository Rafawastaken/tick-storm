package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"

	"github.com/rafawastaken/tick-storm/backend/internal/client/binance"
	"github.com/rafawastaken/tick-storm/backend/internal/crypto"
	"github.com/rafawastaken/tick-storm/backend/internal/ingest"
)

// RunWorker ingests market data until ctx is cancelled or the stream fails.
func RunWorker(ctx context.Context) error {
	b, err := newBase(ctx, "tickstorm-worker")
	if err != nil {
		return err
	}
	defer b.Close()

	svc := crypto.NewService(crypto.NewStore(b.pool))
	w := ingest.NewWorker(
		binance.NewClient(b.cfg.Binance.WSURL, b.log),
		svc,
		b.log,
		b.cfg.Ingest.Symbols,
	)

	g, ctx := errgroup.WithContext(ctx)
	b.serveHTTP(ctx, g, b.httpServer(b.probeRoutes()))
	g.Go(func() error {
		return w.Run(ctx)
	})
	return g.Wait()
}

// probeRoutes is the minimal surface a worker needs: probes today, /metrics next.
func (b *base) probeRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Get("/health", b.handleHealth)
	r.Get("/ready", b.handleReady)

	return r
}
