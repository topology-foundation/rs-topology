// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

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
