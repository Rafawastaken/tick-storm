package ingest

import (
	"fmt"

	"github.com/rafawastaken/tick-storm/backend/internal/client/binance"
	"github.com/rafawastaken/tick-storm/backend/internal/crypto"
	"github.com/shopspring/decimal"
)

const exchangeBinance = "BINANCE"

func toPrice(t binance.Trade) (crypto.Price, error) {
	price, err := decimal.NewFromString(t.Price)
	if err != nil {
		return crypto.Price{}, fmt.Errorf("parse price %q: %w", t.Price, err)
	}
	tradeID := t.TradeID
	return crypto.Price{
		Exchange: exchangeBinance,
		Symbol:   t.Symbol,
		Price:    price,
		TradeID:  &tradeID,
		TradedAt: t.Time(),
	}, nil
}
