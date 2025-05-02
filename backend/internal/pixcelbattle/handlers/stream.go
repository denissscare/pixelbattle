package handlers

import (
	"context"
	"net/http"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WSHandler(svc *service.BattleService, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "upgrader.Upgrade",
				"component": "pixcelbatlle.handlers.WSHandler",
				"success":   false,
			}).Errorf("WS upgrade failed: %s", err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					cancel()
					return
				}
			}
		}()

		canvas, err := svc.InitCanvas(ctx)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "svc.InitCanvas",
				"component": "pixcelbatlle.handlers.WSHandler",
				"success":   false,
			}).Errorf("Init canvas failed: %s", err)
			return
		}
		conn.WriteJSON(map[string]interface{}{"type": "init", "payload": canvas})

		updates, err := svc.Stream(ctx)
		if err != nil {
			log.WithFields(logrus.Fields{
				"action":    "svc.Stream",
				"component": "pixcelbatlle.handlers.WSHandler",
				"success":   false,
			}).Errorf("Stream error: %s", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case px, ok := <-updates:
				if !ok {
					return
				}
				if err := conn.WriteJSON(map[string]interface{}{
					"type":    "update",
					"payload": px,
				}); err != nil {
					_, isClose := err.(*websocket.CloseError)
					if isClose {
						return
					}
					log.WithFields(logrus.Fields{
						"action":    "svc.Stream",
						"component": "pixcelbatlle.handlers.WSHandler",
						"success":   false,
					}).Warnf("write update failed: %s", err)
				}
			}
		}
	}
}
