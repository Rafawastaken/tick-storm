package logger

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware returns a chi middleware that:
//   - extracts or generates W3C trace context
//   - enriches the logger with request_id, trace_id, span_id and injects it into ctx
//   - sets up a per-request fields bag (use logger.Attach to enrich the http_request log)
//   - logs a single "http_request" event on response (including panics, via defer)
//
// Register order: chi.RequestID → logger.Middleware → chi.Recoverer.
// (Logger outside Recoverer so ww.Status() reflects the 500 set by Recoverer.)
func Middleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			tc := traceFromRequest(r)
			requestID := middleware.GetReqID(r.Context())

			l := base.With(
				"request_id", requestID,
				"trace_id", tc.TraceID,
				"span_id", tc.SpanID,
			)

			fields := &ctxFields{}
			ctx := WithLogger(r.Context(), l)
			ctx = WithTrace(ctx, tc)
			ctx = withFields(ctx, fields)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				attrs := []any{
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration_ms", float64(time.Since(start).Microseconds()) / 1000.0,
					"remote_ip", r.RemoteAddr,
				}
				// Append fields attached by downstream handlers/middleware.
				attrs = append(attrs, fields.snapshot()...)

				l.Info("http_request", attrs...)
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

func traceFromRequest(r *http.Request) TraceContext {
	if hdr := r.Header.Get(TraceparentHeader); hdr != "" {
		if tc, err := ParseTraceparent(hdr); err == nil {
			return tc
		}
	}
	return NewTrace()
}

// --- Request-scoped log fields ---

type fieldsCtxKey struct{}

type ctxFields struct {
	mu     sync.Mutex
	fields []any
}

func (f *ctxFields) snapshot() []any {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]any, len(f.fields))
	copy(out, f.fields)
	return out
}

func withFields(ctx context.Context, f *ctxFields) context.Context {
	return context.WithValue(ctx, fieldsCtxKey{}, f)
}

// Attach appends key-value pairs to the http_request log event for the current request.
// Use to enrich the access log with business metadata (user_id, total, search, etc.).
// Pairs follow slog's variadic convention: Attach(ctx, "user_id", 42, "total", 5).
// Safe for concurrent use. No-op if called outside a request handled by Middleware.
func Attach(ctx context.Context, keyvals ...any) {
	if len(keyvals) == 0 {
		return
	}
	f, ok := ctx.Value(fieldsCtxKey{}).(*ctxFields)
	if !ok {
		return
	}
	f.mu.Lock()
	f.fields = append(f.fields, keyvals...)
	f.mu.Unlock()
}
