package log

import (
	"log/slog"
	"os"
)

var defaultOptions = &slog.HandlerOptions{
	Level: slog.LevelInfo,
}

func newSlogLogger(h slog.Handler) *slog.Logger {
	return slog.New(h)
}

func NewSlogTextLogger() *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, defaultOptions)
	return newSlogLogger(ContextHandler{h})
}

func NewSlogJSONLogger() *slog.Logger {
	h := slog.NewJSONHandler(os.Stdout, defaultOptions)
	return newSlogLogger(ContextHandler{h})
}
