package httpx

import "strconv"

// ParseOptFloat returns nil for empty or unparseable input, *float64 otherwise.
// nil means "no filter", which is distinct from zero — a valid filter value.
//
// Note: invalid input is silently treated as absent. Fine for optional filters,
// wrong for anything a caller must be told about.
func ParseOptFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

// ParseOptInt is the int equivalent of ParseOptFloat. Same nil-vs-zero rationale.
func ParseOptInt(s string) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}
