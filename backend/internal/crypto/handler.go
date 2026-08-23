package crypto

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/rafawastaken/tick-storm/backend/pkg/httpx"
)

// Handler adapts HTTP to the service. It owns its own routes so app only
// has to mount them.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{symbol}/latest", h.getLatest)
	r.Get("/{symbol}/prices", h.listPrices)
	return r
}

func (h *Handler) getLatest(w http.ResponseWriter, r *http.Request) {
	price, err := h.svc.GetLatestPrice(r.Context(), GetLatestPriceIn{
		Exchange:   r.URL.Query().Get("exchange"),
		CoinSymbol: chi.URLParam(r, "symbol"),
	})
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(price))
}

func (h *Handler) listPrices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	in := ListPricesForCoinIn{
		Exchange:   q.Get("exchange"),
		CoinSymbol: chi.URLParam(r, "symbol"),
	}
	if limit := httpx.ParseOptInt(q.Get("limit")); limit != nil {
		in.Limit = *limit
	}
	if raw := q.Get("cursor"); raw != "" {
		cursor, err := decodeCursor(raw)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
		in.Before = cursor
	}

	prices, err := h.svc.ListPricesForCoin(r.Context(), in)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	resp := pricesResponse{Data: make([]priceResponse, 0, len(prices))}
	for _, p := range prices {
		resp.Data = append(resp.Data, toResponse(p))
	}
	if len(prices) > 0 {
		resp.NextCursor = encodeCursor(prices[len(prices)-1])
	}

	httpx.WriteJSON(w, http.StatusOK, resp)
}
