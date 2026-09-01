package crypto

import "github.com/go-chi/chi/v5"

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{symbol}/latest", h.getLatest)
	r.Get("/{symbol}/prices", h.listPrices)
	return r
}
