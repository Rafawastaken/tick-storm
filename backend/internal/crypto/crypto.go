package crypto

import (
	"time"

	"github.com/shopspring/decimal"
)

type PriceCursor struct {
	TradedAt time.Time
	ID       int64
}

type Price struct {
	ID       int64
	Exchange string
	Symbol   string
	Price    decimal.Decimal
	TradeID  *int64
	// TradedAt is market time; CreatedAt is ingestion time.
	TradedAt  time.Time
	CreatedAt time.Time
}

type GetLatestPriceIn struct {
	Exchange   string
	CoinSymbol string
}

type ListPricesForCoinIn struct {
	Exchange   string
	CoinSymbol string
	Limit      int
	Before     *PriceCursor
}
