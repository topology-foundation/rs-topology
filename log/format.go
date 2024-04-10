package log

import (
	"bytes"
	"log/slog"
)

const (
	timeFormat        = "2006-01-02T15:04:05"
	termMsgJust       = 40
)

func (h *Handler) format(buf []byte, r slog.Record, usecolor bool) []byte {
	var color = ""
	if usecolor {
		switch r.Level {
		case slog.LevelError:
			color = "\x1b[31m"
		case slog.LevelWarn:
			color = "\x1b[33m"
		case slog.LevelInfo:
			color = "\x1b[32m"
		case slog.LevelDebug:
			color = "\x1b[36m"
		}
	}

	if buf == nil {
		buf = make([]byte, 0, 30+termMsgJust)
	}
	b := bytes.NewBuffer(buf)

	if color != "" { // Start color
		b.WriteString(color)
		b.WriteString(LevelString(r.Level))
		b.WriteString("\x1b[0m")
	} else {
		b.WriteString(LevelString(r.Level))
	}
	b.WriteString("[")
	b.WriteString(r.Time.Format(timeFormat))
	b.WriteString("] ")
	b.WriteString(r.Message)

	h.formatAttributes(b, r, color)
	return b.Bytes()
}

func (h *Handler) formatAttributes(buf *bytes.Buffer, r slog.Record, color string) {
	writeAttr := func(attr slog.Attr, _, _ bool) {
		buf.WriteByte(' ')

		if color != "" {
			buf.WriteString(color)
			buf.WriteString(attr.Key)
			buf.WriteString("\x1b[0m=")
		} else {
			buf.WriteString(attr.Key)
			buf.WriteByte('=')
		}

        buf.WriteString(attr.Value.String())
	}

	var n = 0
	var nAttrs = len(h.attrs) + r.NumAttrs()
	for _, attr := range h.attrs {
		writeAttr(attr, n == 0, n == nAttrs-1)
		n++
	}
	r.Attrs(func(attr slog.Attr) bool {
		writeAttr(attr, n == 0, n == nAttrs-1)
		n++
		return true
	})
	buf.WriteByte('\n')
}
