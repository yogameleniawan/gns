package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/base-go/backend/pkg/config"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	WithRoute(mux *chi.Mux)
	Start() error
	Stop() error
}

type server struct {
	config     *config.Config
	httpServer *http.Server
	mux        *chi.Mux
}

func New() Server {

	return &server{
		config: config.GetConfig(),
	}
}

func (s *server) WithRoute(mux *chi.Mux) {
	s.mux = mux
}

func (s *server) Start() error {
	port := s.config.App.Port
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.mux, // mux from golang chi
	}

	// logging all route
	routes := s.mux.Routes()
	for _, route := range routes {
		log.Println(route.Pattern)
	}

	// run as new goroutine, because we need to start consumer to consume
	// data from message brocker
	go func() {
		log.Printf("Starting server on port %s", port)

		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v\nStack trace:\n%s", r, debug.Stack())
			}
		}()

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	return nil
}

func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.httpServer != nil {
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
