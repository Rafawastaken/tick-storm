package pgxkit

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Helpers to convert between pgx/pgtype values (returned by sqlc) and native
// Go types (used in the domain types).
//
// Convention:
//   - XxxToYyy  — input is Xxx, output is Yyy
//   - Ptr in the name indicates the pointer side (*time.Time)
//   - nil pointer / !Valid pgtype → returns zero-value or nil depending on the target
//
// NOTE: there are deliberately NO float<->Numeric helpers here. A lossy float
// round-trip has no place on price data.

// --- Timestamptz ---

// TimestamptzToTime returns a *time.Time that is nil when t is invalid (NULL in DB).
func TimestamptzToTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}

// TimeToTimestamptz returns an invalid (NULL) Timestamptz when t is nil.
func TimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// --- UUID ---

// UUIDToString returns the canonical 8-4-4-4-12 hyphenated form, or ""
// when u is invalid (NULL in DB).
func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	val, err := u.Value()
	if err != nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// StringToUUID parses a canonical UUID string. Empty string is treated as a
// hard error (not as NULL); use StringPtrToUUID if you want NULL semantics.
func StringToUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}, err
	}
	return u, nil
}

// StringPtrToUUID parses an optional UUID. nil pointer or empty string returns
// an invalid (NULL) pgtype.UUID with no error — for query parameters that
// represent "no value".
func StringPtrToUUID(s *string) (pgtype.UUID, error) {
	if s == nil || *s == "" {
		return pgtype.UUID{}, nil
	}
	return StringToUUID(*s)
}

// UUIDStringPtr is the inverse of StringPtrToUUID — returns nil for an
// invalid UUID, otherwise a pointer to the canonical string.
func UUIDStringPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	s := UUIDToString(u)
	if s == "" {
		return nil
	}
	return &s
}

// --- Primitive pointer shorthand (for sqlc params) ---
//
// sqlc with emit_pointers_for_null_types=true generates *bool / *int /
// *string for nullable columns. These helpers exist so callers can write
// `pgxkit.BoolPtr(true)` instead of `b := true; ...&b` at every call site.

// BoolPtr returns a pointer to v. Shorthand for nullable BOOLEAN columns.
func BoolPtr(v bool) *bool { return &v }

// IntPtr returns a pointer to v. Shorthand for nullable INTEGER columns.
func IntPtr(v int) *int { return &v }

// StringPtr returns a pointer to s, or nil when s is empty. The empty-string
// convention matches "no value" semantics on optional TEXT columns.
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
