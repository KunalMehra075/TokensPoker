package logger

import (
	"log/slog"
	"os"
)

// New builds a structured JSON logger. Production emits JSON; development
// stays JSON too so log shipping is uniform, per the architecture doc.
func New(env string) *slog.Logger {
	level := slog.LevelInfo
	if env != "production" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
