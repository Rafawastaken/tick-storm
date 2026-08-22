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

const (
	serviceName = "tickstorm-api"
	pingTimeout = 5 * time.Second
)

type App struct {
	cfg    *config.Config
	log    *slog.Logger
	pool   *pgxpool.Pool
	server *http.Server
}

// New builds the application and its resources. Fails fast: an unreachable
// database is reported here, not on the first request.
func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.App.Env, logLevel(cfg), serviceName, "")

	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	a := &App{cfg: cfg, log: log, pool: pool}

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: a.routes(),

		// Without these a slow client holds a connection indefinitely.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return a, nil
}

// Run serves HTTP until ctx is cancelled, then drains connections.
// Returns the first error from any component.
func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		a.log.Info("http server listening", "addr", a.server.Addr)
		// A clean shutdown returns ErrServerClosed, which is not a failure.
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		a.log.Info("shutting down", "timeout", a.cfg.App.ShutdownTimeout)

		// A fresh context: ctx is already cancelled and cannot drive the drain.
		shutdownCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx), a.cfg.App.ShutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	})

	return g.Wait()
}

// Close releases resources in reverse order of creation.
func (a *App) Close() {
	if a.pool != nil {
		a.pool.Close()
	}
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
