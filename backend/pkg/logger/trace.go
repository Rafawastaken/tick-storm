package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// W3C Trace Context — https://www.w3.org/TR/trace-context/
// traceparent: 00-<trace_id(32 hex)>-<span_id(16 hex)>-<flags(2 hex)>

// TraceparentHeader is the canonical header name for W3C trace propagation.
const TraceparentHeader = "traceparent"

// TraceContext holds the fields of a W3C traceparent header.
type TraceContext struct {
	TraceID string
	SpanID  string
	Flags   string
}

// NewTrace generates a new root trace with random trace_id and span_id.
// Flags default to "01" (sampled).
func NewTrace() TraceContext {
	return TraceContext{
		TraceID: randHex(16),
		SpanID:  randHex(8),
		Flags:   "01",
	}
}

// NewChildSpan returns a TraceContext with the same TraceID and a fresh SpanID.
// Use this when making outgoing HTTP calls that should belong to the same trace.
func (tc TraceContext) NewChildSpan() TraceContext {
	return TraceContext{
		TraceID: tc.TraceID,
		SpanID:  randHex(8),
		Flags:   tc.Flags,
	}
}

// String returns the traceparent header value for this context.
func (tc TraceContext) String() string {
	return fmt.Sprintf("00-%s-%s-%s", tc.TraceID, tc.SpanID, tc.Flags)
}

// ParseTraceparent parses a W3C traceparent header value.
func ParseTraceparent(s string) (TraceContext, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 4 {
		return TraceContext{}, errors.New("invalid traceparent: expected 4 parts")
	}
	if parts[0] != "00" {
		return TraceContext{}, fmt.Errorf("unsupported traceparent version: %s", parts[0])
	}
	if len(parts[1]) != 32 || len(parts[2]) != 16 || len(parts[3]) != 2 {
		return TraceContext{}, errors.New("invalid traceparent: wrong field lengths")
	}
	return TraceContext{
		TraceID: parts[1],
		SpanID:  parts[2],
		Flags:   parts[3],
	}, nil
}

type traceCtxKey struct{}

// WithTrace returns a new context carrying the given TraceContext.
func WithTrace(ctx context.Context, tc TraceContext) context.Context {
	return context.WithValue(ctx, traceCtxKey{}, tc)
}

// TraceFromContext extracts the TraceContext from ctx.
// Second return is false when no trace is present.
func TraceFromContext(ctx context.Context) (TraceContext, bool) {
	tc, ok := ctx.Value(traceCtxKey{}).(TraceContext)
	return tc, ok
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand should not fail; nothing useful to do
	}
	return hex.EncodeToString(b)
}
