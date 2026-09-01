package ingest

import (
	"context"
	"time"

	"github.com/rafawastaken/tick-storm/backend/internal/client/binance"
	"golang.org/x/sync/errgroup"
)

const reportInterval = 10 * time.Second

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info("ingest worker starting", "symbols", w.symbols)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return w.report(ctx) })
	g.Go(func() error {
		return w.client.StreamTrades(ctx, w.symbols, func(t binance.Trade) error {
			price, err := toPrice(t)
			if err != nil {
				w.skipped.Add(1)
				w.log.Warn("skipping unparseable trade", "trade_id", t.TradeID, "error", err)
				return nil
			}
			if err := w.svc.InsertPrice(ctx, price); err != nil {
				return err
			}
			w.ingested.Add(1)
			w.lagUS.Store(time.Since(t.Time()).Microseconds())
			return nil
		})
	})
	return g.Wait()
}

func (w *Worker) report(ctx context.Context) error {
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			n := w.ingested.Swap(0)
			skipped := w.skipped.Swap(0)
			w.log.Info("ingest progress",
				"trades", n,
				"per_sec", float64(n)/reportInterval.Seconds(),
				"skipped", skipped,
				// Last sample only: a mean or a p99 needs a histogram.
				"lag_ms", float64(w.lagUS.Load())/1000.0,
			)
		}
	}
}
