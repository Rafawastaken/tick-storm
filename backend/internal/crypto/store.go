package crypto

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rafawastaken/tick-storm/backend/internal/crypto/cryptodb"
	"github.com/rafawastaken/tick-storm/backend/pkg/pgxkit"
)

type Store struct {
	pool *pgxpool.Pool
	q    *cryptodb.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool: pool,
		q:    cryptodb.New(pool),
	}
}

func (s *Store) GetLatestPrice(ctx context.Context, in GetLatestPriceIn) (Price, error) {
	row, err := s.q.GetLatestPrice(ctx, cryptodb.GetLatestPriceParams{
		Exchange:   in.Exchange,
		CoinSymbol: in.CoinSymbol,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Price{}, ErrNotFound
		}
		return Price{}, fmt.Errorf("get latest price: %w", err)
	}
	return toDomain(row), nil
}

func (s *Store) InsertPrice(ctx context.Context, p Price) error {
	err := s.q.InsertPrice(ctx, cryptodb.InsertPriceParams{
		Exchange:   p.Exchange,
		CoinSymbol: p.Symbol,
		CoinPrice:  pgxkit.DecimalToNumeric(p.Price),
		TradeID:    p.TradeID,
	})
	if err != nil {
		return fmt.Errorf("insert price: %w", err)
	}
	return nil
}

func (s *Store) ListPricesForCoin(ctx context.Context, in ListPricesForCoinIn) ([]Price, error) {
	params := cryptodb.ListPricesForCoinParams{
		Exchange:   in.Exchange,
		CoinSymbol: in.CoinSymbol,
		MaxResults: int32(in.Limit),
		BeforeTime: pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
		BeforeID:   math.MaxInt64,
	}
	if in.Before != nil {
		params.BeforeTime = pgtype.Timestamptz{Time: in.Before.CreatedAt, Valid: true}
		params.BeforeID = in.Before.ID
	}

	rows, err := s.q.ListPricesForCoin(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list prices for coin: %w", err)
	}
	return toDomainList(rows), nil
}
