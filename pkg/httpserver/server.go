package server

import (
	"context"
	"net/http"
	"time"

	"github.com/yosakoo/task-traker/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.Server.Port,
			Handler:        handler,
			ReadTimeout:    time.Second * 10,
			WriteTimeout:   time.Second * 10,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}