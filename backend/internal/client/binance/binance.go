package binance

import "log/slog"

const (
	defaultBaseURL = "wss://stream.binance.com:9443/stream"
	// The default read limit is 32KB; trade frames are tiny, but a combined
	// stream with many symbols can spike. 1MB is generous and still bounded.
	readLimit = 1 << 20
)

// Client reads public market data. No API key: only account and trading
// endpoints require credentials.
type Client struct {
	baseURL string
	log     *slog.Logger
}

func NewClient(baseURL string, log *slog.Logger) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{baseURL: baseURL, log: log}
}
