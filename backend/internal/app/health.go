package app

import (
	"context"
	"net/http"
	"time"

	"github.com/rafawastaken/tick-storm/backend/pkg/httpx"
)

const readyTimeout = 2 * time.Second

func (b *base) handleHealth(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (b *base) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readyTimeout)
	defer cancel()

	start := time.Now()
	if err := b.pool.Ping(ctx); err != nil {
		httpx.WriteError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	elapsed := time.Since(start)

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":  "ready",
		"ping_ms": float64(elapsed.Microseconds()) / 1000.0,
	})
}
