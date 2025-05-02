package handlers

import (
	"encoding/json"
	"net/http"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"

	"github.com/sirupsen/logrus"
)

func CanvasHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		canvas, err := svc.InitCanvas(r.Context())
		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "svc.InitCanvas",
				"component": "pixcelbatlle.handlers.CanvasHandler",
				"success":   false,
			}).Errorf("Init canvas failed: %s", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(canvas)
	}
}
