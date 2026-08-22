package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/rafawastaken/tick-storm/backend/pkg/logger"
)

// ErrorResponse is the standard error envelope emitted by WriteError.
// Exported so API docs (Swagger/OpenAPI) can reference it as a schema.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes body as JSON with the given status.
// Ignores encode errors silently — the response is already committed.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteError writes a standard error envelope: {"error": "<msg>"}.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}

// WriteInternal attaches err to the request's log context and writes the
// generic 500 envelope. Single place an unclassified error is recorded, so an
// error that reaches here is logged exactly once and never leaks to the client.
func WriteInternal(w http.ResponseWriter, r *http.Request, err error) {
	logger.Attach(r.Context(), "error", err.Error())
	WriteError(w, http.StatusInternalServerError, "internal error")
}
