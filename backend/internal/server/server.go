package server

import (
	"fmt"
	"net/http"
	"pixelbattle/internal/config"

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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("failed to start server: " + err.Error())
		}
	}()

	return nil
}

func InitRouter() *chi.Mux {
	router := chi.NewRouter()
	//TODO: add middlewares and hadlers for routes
	return router
}
