package crypto

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rafawastaken/tick-storm/backend/internal/crypto/cryptodb"
)

// Store is the only place that talks to the database. It translates pgtype
// values into domain types so no other layer imports pgx.
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
	if err := s.q.InsertPrice(ctx, toInsertParams(p)); err != nil {
		return fmt.Errorf("insert price: %w", err)
	}
	return nil
}
