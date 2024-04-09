package log

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

type emptyHandler struct{}

// DiscardHandler returns a no-op handler
func EmptyHandler() slog.Handler {
	return &emptyHandler{}
}

func (h *emptyHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (h *emptyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return false
}

func (h *emptyHandler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func (h *emptyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &emptyHandler{}
}

type Handler struct {
	mu           sync.Mutex
	wr           io.Writer
	lvl          slog.Level
	useColor     bool
	attrs        []slog.Attr
	fieldPadding map[string]int

	buf []byte
}

func NewHandler(wr io.Writer, lvl slog.Level, useColor bool) *Handler {
	return &Handler{
		wr:           wr,
		lvl:          lvl,
		useColor:     useColor,
		fieldPadding: make(map[string]int),
	}
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	buf := h.format(h.buf, r, h.useColor)
	h.wr.Write(buf)
	h.buf = buf[:0]
	return nil
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.lvl
}

func (h *Handler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("not implemented")
}
