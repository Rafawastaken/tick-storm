package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/coder/websocket"
)

// StreamTrades connects and calls fn for every trade until ctx is cancelled
// or the connection drops. Returning an error is the caller's cue to reconnect.
func (c *Client) StreamTrades(ctx context.Context, symbols []string, fn func(Trade) error) error {
	url, err := c.streamURL(symbols)
	if err != nil {
		return err
	}

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("dial binance: %w", err)
	}
	defer conn.CloseNow()
	conn.SetReadLimit(readLimit)

	c.log.Info("binance stream connected", "symbols", symbols)

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("read binance stream: %w", err)
		}

		var msg streamMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.log.Warn("malformed envelope, skipping", "error", err)
			continue
		}

		var trade Trade
		if err := json.Unmarshal(msg.Data, &trade); err != nil {
			c.log.Warn("malformed trade, skipping", "stream", msg.Stream, "error", err)
			continue
		}
		if trade.EventType != "trade" {
			continue
		}

		if err := fn(trade); err != nil {
			return err
		}
	}
}

func (c *Client) streamURL(symbols []string) (string, error) {
	if len(symbols) == 0 {
		return "", errors.New("no symbols to subscribe")
	}
	streams := make([]string, 0, len(symbols))
	for _, s := range symbols {
		streams = append(streams, strings.ToLower(strings.TrimSpace(s))+"@trade")
	}
	return c.baseURL + "?streams=" + strings.Join(streams, "/"), nil
}
