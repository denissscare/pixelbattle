package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"

	"github.com/sirupsen/logrus"
)

func UpdatePixelHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p domain.Pixel
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			log.WithFields(logrus.Fields{
				"action":    "json.NewDecoder",
				"component": "pixcelbatlle.handlers.UpdatePixelHandler",
				"success":   false,
			}).Errorf("bad request payload: %s", err)
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		if err := svc.UpdatePixel(r.Context(), p); err != nil {
			log.WithFields(logrus.Fields{
				"action":    "svc.UpdatePixel",
				"component": "pixcelbatlle.handlers.UpdatePixelHandler",
				"success":   false,
			}).Errorf("bad request payload: %s", err)
			http.Error(w, "failed to update pixel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
