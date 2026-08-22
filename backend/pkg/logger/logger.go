package logger

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/lmittmann/tint"
)

// New returns a *slog.Logger configured for the given environment, with
// service identity attributes attached as baseline (`service`,
// `service_version`, `env`) so every log line is self-describing when
// ingested by Grafana/Loki/Datadog/Elastic — no extra mapping rules required.
//
// env "prod" emits JSON; anything else emits colored text via tint.
// Pass version="" to derive it from build info (`vcs.revision` short SHA when
// available, "dev" otherwise — typical under `go run`).
func New(env string, level slog.Level, service, version string) *slog.Logger {
	if version == "" {
		version = buildVersion()
	}

	var handler slog.Handler
	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: "15:04:05.000",
		})
	}

	return slog.New(handler).With(
		"service", service,
		"service_version", version,
		"env", env,
	)
}

// buildVersion returns the VCS short revision from build info, or "dev" when
// not embedded (e.g. under `go run` without -buildvcs).
func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && len(s.Value) >= 7 {
			return s.Value[:7]
		}
	}
	return "dev"
}
