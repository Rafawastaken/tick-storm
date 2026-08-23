package crypto

import (
	"time"

	"github.com/shopspring/decimal"
)

type PriceCursor struct {
	CreatedAt time.Time
	ID        int64
}

type Price struct {
	ID        int64
	Exchange  string
	Symbol    string
	Price     decimal.Decimal
	TradeID   *int64
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
