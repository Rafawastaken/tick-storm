package binance

import (
	"context"
	"math/rand/v2"
	"time"
)

const (
	initialBackoff   = 1 * time.Second
	maxBackoff       = 30 * time.Second
	stableConnection = 1 * time.Minute
)

func (c *Client) StreamTradesWithRetry(ctx context.Context, symbols []string, fn func(Trade) error) error {
	backoff := initialBackoff

	for {
		start := time.Now()
		err := c.StreamTrades(ctx, symbols, fn)
		if ctx.Err() != nil {
			return nil
		}

		if time.Since(start) >= stableConnection {
			backoff = initialBackoff
		}

		wait := withJitter(backoff)
		c.log.Warn("binance stream lost, reconnecting",
			"error", err,
			"connected_for", time.Since(start).Round(time.Second),
			"retry_in", wait.Round(time.Millisecond),
		)

		if err := sleepCtx(ctx, wait); err != nil {
			return nil
		}

		backoff = min(backoff*2, maxBackoff)
	}
}

// sleepCtx waits for d, or returns early when ctx is cancelled.
func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// withJitter spreads reconnects so many clients do not retry in lockstep.
func withJitter(d time.Duration) time.Duration {
	return d/2 + time.Duration(rand.Int64N(int64(d/2)))
}
