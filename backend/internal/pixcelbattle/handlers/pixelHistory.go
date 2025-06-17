package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
)

func PixelHistoryHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		history, err := svc.GetAllPixelHistory(r.Context())
		if err != nil {
			log.Errorf("HTTP GET /pixels/history failed: %v", err)
			http.Error(w, "failed to get pixel history", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(history); err != nil {
			log.Errorf("HTTP GET /pixels/history encode failed: %v", err)
			http.Error(w, "failed to encode pixel history", http.StatusInternalServerError)
			return
		}

		log.Infof("HTTP GET /pixels/history → returned %d records", len(history))
	}
}
