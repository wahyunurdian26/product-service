package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (s *HttpServer) Run() {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		log.Printf("Server Product started successfully - Running on port :%s", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	s.waitForShutdown(srv)
}

func (s *HttpServer) waitForShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")
	s.Stop(srv)
}

func (s *HttpServer) Stop(srv *http.Server) {
	log.Println("Stopping Server - initiating graceful shutdown")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	// Close database & redis connections
	if s.close != nil {
		s.close()
	}

	log.Println("Server stopped successfully")
}
