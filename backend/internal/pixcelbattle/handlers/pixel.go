package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
)

func UpdatePixelHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p domain.Pixel
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			log.Errorf("HTTP POST /pixel bad payload: %v", err)
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}
		if err := svc.UpdatePixel(r.Context(), p); err != nil {
			log.Errorf("HTTP POST /pixel update failed: %v", err)
			http.Error(w, "failed to update pixel", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.Infof("HTTP POST /pixel → updated pixel at (%d,%d) by %s", p.X, p.Y, p.Author)
	}
}
