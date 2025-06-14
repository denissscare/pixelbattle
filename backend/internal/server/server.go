package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pixelbattle/internal/config"
	"pixelbattle/internal/middleware"
	"pixelbattle/internal/pixcelbattle/handlers"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg    *config.Config
	router *chi.Mux
	log    *logger.Logger
}

func New(cfg *config.Config, router *chi.Mux, log *logger.Logger) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
		log:    log,
	}
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%d", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler:      s.router,
		ReadTimeout:  s.cfg.Server.Timeout,
		WriteTimeout: s.cfg.Server.Timeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	go func() {
		s.log.WithFields(logrus.Fields{
			"action":    "server.Run",
			"component": "internal.server.Run",
			"success":   true,
		}).Infof("server starts on: %v:%d", s.cfg.Server.Host, s.cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.WithFields(logrus.Fields{
				"action":    "server.ListenAndServe",
				"component": "internal.server.Run",
				"success":   false,
			}).Errorf("Failed to start server: %v", err)

		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.log.WithFields(logrus.Fields{
		"action":    "server.getChannel",
		"component": "internal.server.Run",
		"success":   true,
	}).Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.log.WithFields(logrus.Fields{
			"action":    "server.Shutdown",
			"component": "internal.server.Run",
			"success":   false,
		}).Errorf("Failed to gracefully stop the server: %v", err)
	}

	s.log.WithFields(logrus.Fields{
		"action":    "server.Shutdown",
		"component": "internal.server.Run",
		"success":   true,
	}).Info("Server stopped successfully")

	return nil
}

func InitRouter(svc *service.BattleService, log *logger.Logger, metrics metrics.Metrics) *chi.Mux {
	router := chi.NewRouter()

	router.With(middleware.NoLogger).Get("/ws", handlers.WSHandler(svc, log))

	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/canvas", handlers.CanvasHandler(svc, log))
	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/pixel", handlers.UpdatePixelHandler(svc, log))

	router.Handle("/metrics", promhttp.Handler())

	return router
}
