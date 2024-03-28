package logs

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

type Handler struct {
	slog.Handler
	opts    slog.HandlerOptions
	w       io.Writer
	m       *sync.Mutex
	tracker *progressTracker
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if h.tracker != nil {
		h.tracker.msg <- r.Message
		if r.Level == slog.LevelInfo {
			h.m.Lock()
			h.tracker.Increment(1)
			h.m.Unlock()
		}
	}
	return h.Handler.Handle(ctx, r)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler := &Handler{}
	*handler = *h
	handler.Handler = h.Handler.WithAttrs(attrs)
	return handler
}

func (h *Handler) WithGroup(name string) slog.Handler {
	handler := &Handler{}
	*handler = *h
	handler.Handler = h.Handler.WithGroup(name)
	return handler
}

func NewHandler(
	out io.Writer,
	opts slog.HandlerOptions,
) *Handler {
	return &Handler{
		opts:    opts,
		w:       out,
		m:       &sync.Mutex{},
		Handler: slog.NewJSONHandler(out, &opts),
	}
}
