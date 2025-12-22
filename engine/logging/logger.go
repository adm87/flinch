package logging

import (
	"io"
	"log/slog"
)

// Logger is a simple logging interface.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	With(args ...any) Logger
}

type logger struct {
	slog *slog.Logger
}

// NewLogger creates a new Logger instance with a slog backend.
func NewLogger(writer io.Writer, args ...any) Logger {
	return NewLoggerWithOptions(writer, &slog.HandlerOptions{}, args...)
}

// NewLoggerWithOptions creates a new Logger instance with a slog backend and custom handler options.
func NewLoggerWithOptions(writer io.Writer, options *slog.HandlerOptions, args ...any) Logger {
	return &logger{slog: slog.New(slog.NewTextHandler(writer, options)).With(args...)}
}

func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *logger) With(args ...any) Logger {
	return &logger{
		slog: l.slog.With(args...),
	}
}
