package crypto

import (
	"errors"
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
	return r
}

// writeError maps domain errors to status codes. Anything unmapped is a bug
// or an outage, so it goes to WriteInternal and gets logged there.
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
