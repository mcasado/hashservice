package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
	"github.com/mcasado/hashservice/stats"
)


var (
	healthy         int32
	quit            = make(chan os.Signal, 1)
	statsMiddleware = stats.New()
)

// Server implements our HTTP server
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP Server
func newServer(port string, h http.Handler, l *log.Logger) *Server {
	handler := statsMiddleware.Handler(h)
	return &Server{
		server: &http.Server{
			Addr:           ":" + port,
			Handler:        handler, // pass in mux/router
			ErrorLog:       l,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

// Run starts the HTTP server
func (s *Server) run() error {

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	done := make(chan bool)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-quit
		fmt.Println("")
		s.server.ErrorLog.Printf("%s - Shutdown signal received...\n", hostname)
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)
		if err := s.server.Shutdown(ctx); err != nil {
			s.server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	s.server.ErrorLog.Printf("%s - Starting server on port %v", hostname, s.server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.server.ErrorLog.Fatalf("Could not listen on %s: %v\n", s.server.Addr, err)
	}

	<-done
	time.Sleep(5 * time.Second)
	s.server.ErrorLog.Printf("%s - Server gracefully stopped.\n", hostname)
	return nil
}