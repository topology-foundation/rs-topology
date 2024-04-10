package log

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

const errorKey = "LOG_ERROR"

// LevelString returns a string containing the name of a Lvl.
func LevelString(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return "unknown"
	}
}

func StringLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		// defaults to info
		return slog.LevelInfo
	}
}

type Logger interface {
	With(ctx ...interface{}) Logger

	New(ctx ...interface{}) Logger

	Log(level slog.Level, msg string, ctx ...interface{})

	Debug(msg string, ctx ...interface{})

	Info(msg string, ctx ...interface{})

	Warn(msg string, ctx ...interface{})

	Error(msg string, ctx ...interface{})

	Write(level slog.Level, msg string, attrs ...any)

	Handler() slog.Handler
}

type logger struct {
	inner *slog.Logger
}

// NewLogger returns a logger with the specified handler set
func NewLogger(h slog.Handler) Logger {
	return &logger{
		slog.New(h),
	}
}

func (l *logger) Handler() slog.Handler {
	return l.inner.Handler()
}

// Write logs a message at the specified level:
func (l *logger) Write(level slog.Level, msg string, attrs ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	if len(attrs)%2 != 0 {
		attrs = append(attrs, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(attrs...)
	err := l.inner.Handler().Handle(context.Background(), r)
	if err != nil {
		Error("(Logger) Error writing a message", "error", err)
	}
}

func (l *logger) Log(level slog.Level, msg string, attrs ...any) {
	l.Write(level, msg, attrs...)
}

func (l *logger) With(ctx ...interface{}) Logger {
	return &logger{l.inner.With(ctx...)}
}

func (l *logger) New(ctx ...interface{}) Logger {
	return l.With(ctx...)
}

func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.Write(slog.LevelDebug, msg, ctx...)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	l.Write(slog.LevelInfo, msg, ctx...)
}

func (l *logger) Warn(msg string, ctx ...any) {
	l.Write(slog.LevelWarn, msg, ctx...)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	l.Write(slog.LevelError, msg, ctx...)
}
