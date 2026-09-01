package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/rafawastaken/tick-storm/backend/internal/config"
	"github.com/rafawastaken/tick-storm/backend/pkg/logger"
)

const pingTimeout = 5 * time.Second

// base holds what every process needs, whatever its role.
type base struct {
	cfg  *config.Config
	log  *slog.Logger
	pool *pgxpool.Pool
}

// newBase loads configuration and opens shared resources. Fails fast: an
// unreachable database is reported here, not on the first request.
func newBase(ctx context.Context, service string) (*base, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.App.Env, logLevel(cfg), service, "")

	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &base{cfg: cfg, log: log, pool: pool}, nil
}

func (b *base) Close() {
	if b.pool != nil {
		b.pool.Close()
	}
}

func (b *base) httpServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", b.cfg.App.Port),
		Handler: h,

		// Without these a slow client holds a connection indefinitely.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// serveHTTP adds the server and its shutdown to g.
func (b *base) serveHTTP(ctx context.Context, g *errgroup.Group, srv *http.Server) {
	g.Go(func() error {
		b.log.Info("http server listening", "addr", srv.Addr)
		// A clean shutdown returns ErrServerClosed, which is not a failure.
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		b.log.Info("shutting down", "timeout", b.cfg.App.ShutdownTimeout)

		// A fresh context: ctx is already cancelled and cannot drive the drain.
		shutdownCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx), b.cfg.App.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	})
}

func newPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}
	poolCfg.MaxConns = cfg.Postgres.MaxConns
	poolCfg.MinConns = cfg.Postgres.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	// Bounded: without a timeout an unreachable host stalls startup.
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

func logLevel(cfg *config.Config) slog.Level {
	if cfg.App.IsProduction() {
		return slog.LevelInfo
	}
	return slog.LevelDebug
}
