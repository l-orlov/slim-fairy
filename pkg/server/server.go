package server

import (
	"context"
	"net/http"
	"time"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
	timeout        = 10 * time.Second
)

type (
	Server struct {
		httpServer *http.Server
	}
)

func New(handler http.Handler, addres string) *Server {
	s := &Server{}
	s.httpServer = &http.Server{
		Addr:           addres,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
	}

	return s
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
