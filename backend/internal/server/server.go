package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pixelbattle/internal/config"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	cfg    *config.Config
	router *chi.Mux
}

func New(cfg *config.Config, router *chi.Mux) *Server {
	return &Server{
		cfg:    cfg,
		router: router,
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
		fmt.Printf("\n\nSERVER START ON: %v:%d\n\n", s.cfg.Server.Host, s.cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("failed to start server: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Printf("\n\nSERVER STOPPING\n\n")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic("failed to stop server")
	}
	fmt.Printf("\nSERVER STOP\n")

	return nil
}

func InitRouter() *chi.Mux {
	router := chi.NewRouter()
	//TODO: add middlewares and hadlers for routes
	return router
}
