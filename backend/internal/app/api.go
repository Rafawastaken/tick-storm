package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"

	"github.com/rafawastaken/tick-storm/backend/internal/crypto"
	"github.com/rafawastaken/tick-storm/backend/pkg/logger"
)

// RunAPI serves the HTTP API until ctx is cancelled.
func RunAPI(ctx context.Context) error {
	b, err := newBase(ctx, "tickstorm-api")
	if err != nil {
		return err
	}
	defer b.Close()

	svc := crypto.NewService(crypto.NewStore(b.pool))
	srv := b.httpServer(b.apiRoutes(crypto.NewHandler(svc)))

	g, ctx := errgroup.WithContext(ctx)
	b.serveHTTP(ctx, g, srv)
	return g.Wait()
}

// apiRoutes builds the public API. Slice routes get mounted here.
func (b *base) apiRoutes(cryptoHandler *crypto.Handler) http.Handler {
	const apiPrefix = "/api/v1"

	r := chi.NewRouter()

	// Order matters: RequestID feeds the logger, and the logger sits outside
	// Recoverer so a panic is still logged with its 500 status.
	r.Use(middleware.RequestID)
	r.Use(logger.Middleware(b.log))
	r.Use(middleware.Recoverer)

	// Probes stay unversioned: they are infrastructure, not API.
	r.Get("/health", b.handleHealth)
	r.Get("/ready", b.handleReady)

	r.Mount(apiPrefix+"/crypto", cryptoHandler.Routes())

	return r
}
