package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
	"time"
)

type PixelRequest struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Color  string `json:"color"`
	Author string `json:"author"`
}

func UpdatePixelHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PixelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("HTTP POST /pixel bad payload: %v", err)
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		var p domain.Pixel = domain.Pixel{
			X:         req.X,
			Y:         req.Y,
			Color:     req.Color,
			Author:    req.Author,
			Timestamp: time.Now(),
		}
		if err := svc.UpdatePixel(r.Context(), p); err != nil {
			if verrs, ok := err.(domain.ValidationErrors); ok {
				log.Infof("HTTP POST /pixel update failed: %v", verrs)
				http.Error(w, verrs.Error(), http.StatusBadRequest)
				return
			}
			log.Errorf("HTTP POST /pixel update failed: %v", err)
			http.Error(w, "failed to update pixel", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		log.Infof("HTTP POST /pixel → updated pixel at (%d,%d) by %s", p.X, p.Y, p.Author)
	}
}
