package handlers

import (
	"log/slog"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("PONG")); err != nil {
		slog.Error("Ping message send failed: ", "error", err.Error())
		return
	}
}
