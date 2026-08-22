package crypto

import (
	"github.com/rafawastaken/tick-storm/backend/internal/crypto/cryptodb"
	"github.com/rafawastaken/tick-storm/backend/pkg/pgxkit"
)

func toDomain(row cryptodb.CryptoPrice) Price {
	return Price{
		ID:        row.ID,
		Exchange:  row.Exchange,
		Symbol:    row.CoinSymbol,
		Price:     pgxkit.NumericToDecimal(row.CoinPrice),
		TradeID:   row.TradeID,
		CreatedAt: row.CreatedAt.Time,
	}
}

func toInsertParams(p Price) cryptodb.InsertPriceParams {
	return cryptodb.InsertPriceParams{
		Exchange:   p.Exchange,
		CoinSymbol: p.Symbol,
		CoinPrice:  pgxkit.DecimalToNumeric(p.Price),
		TradeID:    p.TradeID,
	}
}
