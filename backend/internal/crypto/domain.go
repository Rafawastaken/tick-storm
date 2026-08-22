package crypto

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// Domain errors. Callers branch on these instead of on storage details, so
// nothing above this file needs to know pgx exists.
var (
	ErrNotFound      = errors.New("price not found")
	ErrInvalidSymbol = errors.New("invalid symbol")
)

// Price is a single observation of a coin's price.
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
