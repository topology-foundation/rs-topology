// The code inside the `log` directory is a modified version of the [`go-ethereum` log](https://github.com/ethereum/go-ethereum/tree/bd91810462187086b2715fd343aa427e181d89a2/log)

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
	mutex    sync.Mutex
	writer   io.Writer
	lvl      slog.Level
	useColor bool
	buf      []byte
}

func NewHandler(writer io.Writer, lvl slog.Level, useColor bool) *Handler {
	return &Handler{
		writer:   writer,
		lvl:      lvl,
		useColor: useColor,
	}
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	buf := h.format(h.buf, record, h.useColor)
	_, err := h.writer.Write(buf)

	if err != nil {
		return err
	}

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
