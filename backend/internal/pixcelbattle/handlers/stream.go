package handlers

import (
	"context"
	"net/http"
	"pixelbattle/internal/pixcelbattle/domain"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WSHandler(svc *service.BattleService, log *logger.Logger, timeout int) http.HandlerFunc {
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

		author := r.URL.Query().Get("username")
		var lastPixel *domain.Pixel
		if author != "" {
			lastPixel, err = svc.GetLastPixelByAuthor(ctx, author)
			if err != nil {
				log.Warnf("WS: could not fetch last pixel for %s: %v", author, err)
			}
		}
		conn.WriteJSON(map[string]interface{}{
			"type": "init",
			"payload": map[string]interface{}{
				"timeout":    timeout,
				"canvas":     canvas,
				"last_pixel": lastPixel,
			},
		})
		log.Infof("WS: sent initial canvas (%d pixels) + last pixel to %s", len(canvas), r.RemoteAddr)

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
				if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
					log.Warnf("WS: set deadline failed for %s: %v", r.RemoteAddr, err)
				}
				if err := conn.WriteJSON(map[string]interface{}{
					"type":    "update",
					"payload": px,
				}); err != nil {
					_, isClose := err.(*websocket.CloseError)
					if isClose {
						return
					}
					log.Warnf("WS: write update failed for %s: %v", r.RemoteAddr, err)
				} else {
					log.Debugf("WS: sent update %+v to %s", px, r.RemoteAddr)
				}
			}
		}
	}
}
