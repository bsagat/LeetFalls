package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(srv *http.Server) {
	go func() {
		slog.Info("Server has been started on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %s", err.Error())
		}
	}()
}

// Listens for system signals (e.g., SIGINT, SIGTERM) to ensure a graceful shutdown of the HTTP server.
func WaitForShutDown(srv *http.Server) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	slog.Info("ShutDown signal received!!!")
	slog.Info("Shutting down HTTP server")

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(context); err != nil {
		slog.Error("Server shutdown failed: ", "error", err.Error())
		return
	}

	slog.Info("HTTP server gracefully stopped")
}
