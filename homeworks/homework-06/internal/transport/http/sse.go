package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"homework06/internal/requestid"
)

func (h *handler) events(w http.ResponseWriter, r *http.Request) {
	// Контекст подписки связан с HTTP-запросом, но не с unary-таймаутом.
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	stream, err := h.client.Watch(ctx, r.PathValue("id"))
	if err != nil {
		fail(w, r, err)
		return
	}
	if stream == nil {
		fail(w, r, errors.New("nil backend stream"))
		return
	}
	first, err := stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("stream ended before first snapshot")
		}
		fail(w, r, err)
		return
	}
	if first == nil {
		fail(w, r, errors.New("nil stream response"))
		return
	}
	firstData, err := json.Marshal(normalize(first))
	if err != nil {
		fail(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	controller := http.NewResponseController(w)
	send := func(kind string, data []byte) error {
		// Ограничиваем время записи медленному HTTP-клиенту.
		if err := controller.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil && !errors.Is(err, http.ErrNotSupported) {
			return err
		}
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", kind, data); err != nil {
			return err
		}
		return controller.Flush()
	}
	defer func() {
		if err := controller.SetWriteDeadline(time.Time{}); err != nil && !errors.Is(err, http.ErrNotSupported) {
			slog.Debug("reset write deadline", "error", err)
		}
	}()
	if err := send("snapshot", firstData); err != nil {
		return
	}
	for {
		state, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err == nil && state == nil {
			err = errors.New("nil stream response")
		}
		var data []byte
		if err == nil {
			data, err = json.Marshal(normalize(state))
		}
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			_, message := publicError(err)
			slog.Error("subscription failed", "request_id", requestid.FromContext(ctx), "error", err)
			errorData, _ := json.Marshal(map[string]string{"error": message, "request_id": requestid.FromContext(ctx)})
			if sendErr := send("error", errorData); sendErr != nil {
				slog.Debug("send SSE error", "error", sendErr)
			}
			return
		}
		if err := send("snapshot", data); err != nil {
			return
		}
	}
}
