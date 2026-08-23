package crypto

import "errors"

// Domain errors. Callers branch on these instead of on storage details, so
// nothing above this file needs to know pgx exists.
var (
	ErrNotFound      = errors.New("price not found")
	ErrInvalidSymbol = errors.New("invalid symbol")
)
