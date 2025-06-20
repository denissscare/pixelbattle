package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	authHandlers "pixelbattle/internal/auth/handlers"
	auth "pixelbattle/internal/auth/service"
	"pixelbattle/internal/config"
	"pixelbattle/internal/middleware"
	"pixelbattle/internal/pixcelbattle/handlers"
	"pixelbattle/internal/pixcelbattle/metrics"
	"pixelbattle/internal/pixcelbattle/service"
	"pixelbattle/internal/s3"
	jwtutil "pixelbattle/pkg/jwt"
	"pixelbattle/pkg/logger"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

func InitRouter(svc *service.BattleService,
	authSvc *auth.Service,
	log *logger.Logger,
	metrics metrics.Metrics,
	jwtManager *jwtutil.JWTManager,
	s3Client *s3.Client,
	cfg config.Config) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           int((300 * time.Second).Seconds()),
	}))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})

	router.With(middleware.NoLogger).Get("/ws", handlers.WSHandler(svc, log, cfg.Server.CanvasTimeout))

	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/canvas", handlers.CanvasHandler(svc, log))
	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/pixel", handlers.UpdatePixelHandler(svc, log))
	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/index", handlers.CanvasRender(svc, log))
	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/register", authHandlers.RegisterHandler(s3Client, authSvc, log))
	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/register", authHandlers.RegisterRender(authSvc, log))
	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/login", authHandlers.LoginRender(authSvc, log))
	router.With(middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/login", authHandlers.LoginHandler(authSvc, log, cfg.Minio.PublicHost))
	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/avatar", authHandlers.UploadAvatarHandler(authSvc, log))
	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Post("/email", authHandlers.UpdateEmailHandler(authSvc, log))

	router.With(middleware.JWTAuth(jwtManager), middleware.Metrics(metrics), middleware.RequestLogger(log)).Get("/pixels/history", handlers.PixelHistoryHandler(svc, log))

	router.Handle("/metrics", promhttp.Handler())

	return router
}