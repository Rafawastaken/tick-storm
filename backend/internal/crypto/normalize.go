package crypto

import (
	"strings"

	"github.com/rafawastaken/tick-storm/backend/internal/crypto/cryptodb"
	"github.com/rafawastaken/tick-storm/backend/pkg/pgxkit"
)

const (
	defaultExchange = "BINANCE"
	defaultLimit    = 100
	maxLimit        = 1000
)

func toDomain(row cryptodb.CryptoPrice) Price {
	return Price{
		ID:        row.ID,
		Exchange:  row.Exchange,
		Symbol:    row.CoinSymbol,
		Price:     pgxkit.NumericToDecimal(row.CoinPrice),
		TradeID:   row.TradeID,
		TradedAt:  row.TradedAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}
}

func toDomainList(rows []cryptodb.CryptoPrice) []Price {
	prices := make([]Price, 0, len(rows))
	for _, row := range rows {
		prices = append(prices, toDomain(row))
	}
	return prices
}

func normalizeExchange(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return defaultExchange
	}
	return s
}

func normalizeCoinSymbol(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func clampLimit(n int) int {
	if n <= 0 {
		return defaultLimit
	}
	if n > maxLimit {
		return maxLimit
	}
	return n
}
