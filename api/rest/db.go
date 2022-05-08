package rest

import (
	"context"
	"net/http"
	"time"
)

func (h *Handler) heartbeat(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
	defer cancel()

	if err := h.Server.Heartbeat(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
