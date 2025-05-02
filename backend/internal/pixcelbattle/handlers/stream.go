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
		log.Infof("WS: new connection from %s", r.RemoteAddr)
		defer func() {
			conn.Close()
			log.Infof("WS: client %s disconnected", r.RemoteAddr)
		}()

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
			log.Errorf("WS: InitCanvas error: %v", err)
			return
		}
		conn.WriteJSON(map[string]interface{}{"type": "init", "payload": canvas})
		log.Infof("WS: sent initial canvas (%d pixels) to %s", len(canvas), r.RemoteAddr)

		updates, err := svc.Stream(ctx)
		if err != nil {
			log.Errorf("WS: Stream subscribe error: %v", err)
			return
		}
		log.Infof("WS: subscribed to canvas.updates for %s", r.RemoteAddr)

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
					}).Warnf("WS: write update failed for %s: %v", r.RemoteAddr, err)
				} else {
					log.Debugf("WS: sent update %+v to %s", px, r.RemoteAddr)
				}
			}
		}
	}
}
