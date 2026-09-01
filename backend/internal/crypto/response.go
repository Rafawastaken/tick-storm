package crypto

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rafawastaken/tick-storm/backend/pkg/httpx"
)

type priceResponse struct {
	Exchange string    `json:"exchange"`
	Symbol   string    `json:"symbol"`
	Price    string    `json:"price"`
	TradedAt time.Time `json:"traded_at"`
}

type pricesResponse struct {
	Data       []priceResponse `json:"data"`
	NextCursor string          `json:"next_cursor,omitempty"`
}

func toResponse(p Price) priceResponse {
	return priceResponse{
		Exchange: p.Exchange,
		Symbol:   p.Symbol,
		Price:    p.Price.String(),
		TradedAt: p.TradedAt,
	}
}

func encodeCursor(p Price) string {
	raw := fmt.Sprintf("%d:%d", p.TradedAt.UnixNano(), p.ID)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(s string) (*PriceCursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	nanos, id, found := strings.Cut(string(raw), ":")
	if !found {
		return nil, errors.New("malformed cursor")
	}
	n, err := strconv.ParseInt(nanos, 10, 64)
	if err != nil {
		return nil, err
	}
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return &PriceCursor{TradedAt: time.Unix(0, n), ID: i}, nil
}

func (h *Handler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "price not found")
	case errors.Is(err, ErrInvalidSymbol):
		httpx.WriteError(w, http.StatusBadRequest, "invalid symbol")
	default:
		httpx.WriteInternal(w, r, err)
	}
}
