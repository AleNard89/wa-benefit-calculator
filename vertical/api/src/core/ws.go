package core

import (
	"net/http"

	"go.uber.org/zap"
)

func HandleWs(w http.ResponseWriter, r *http.Request) {
	// WebSocket hub placeholder - will be implemented in Phase 2
	zap.S().Debug("WebSocket connection attempted")
	w.WriteHeader(http.StatusNotImplemented)
}
