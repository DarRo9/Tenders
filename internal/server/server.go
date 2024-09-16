package server

import (
	"context"
	"net/http"
	"time"

	"github.com/DarRo9/Service_for_creating_and_processing_tenders/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func New(handler http.Handler, cfg *config.ServerConfig) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.Address,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
