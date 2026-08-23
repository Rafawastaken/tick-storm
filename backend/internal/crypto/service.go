package crypto

import (
	"context"
)

// Service holds the business rules. It knows the store, never the transport.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetLatestPrice(ctx context.Context, in GetLatestPriceIn) (Price, error) {
	in.Exchange = normalizeExchange(in.Exchange)
	in.CoinSymbol = normalizeCoinSymbol(in.CoinSymbol)
	if in.CoinSymbol == "" {
		return Price{}, ErrInvalidSymbol
	}
	return s.store.GetLatestPrice(ctx, in)
}

func (s *Service) ListPricesForCoin(ctx context.Context, in ListPricesForCoinIn) ([]Price, error) {
	in.Exchange = normalizeExchange(in.Exchange)
	in.CoinSymbol = normalizeCoinSymbol(in.CoinSymbol)
	if in.CoinSymbol == "" {
		return nil, ErrInvalidSymbol
	}
	in.Limit = clampLimit(in.Limit)
	return s.store.ListPricesForCoin(ctx, in)
}

func (s *Service) InsertPrice(ctx context.Context, p Price) error {
	p.Exchange = normalizeExchange(p.Exchange)
	p.Symbol = normalizeCoinSymbol(p.Symbol)
	if p.Symbol == "" {
		return ErrInvalidSymbol
	}
	return s.store.InsertPrice(ctx, p)
}
