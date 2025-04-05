package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taldoflemis/nume/configs"
	"github.com/taldoflemis/nume/internal/server"
)

func gracefulShutdown(
	apiServer *http.Server,
	done chan bool,
	shutdownTimeoutInSeconds int,
) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully. press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(shutdownTimeoutInSeconds)*time.Second,
	)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
		slog.Error("server forced to shutdown", slog.Any("error", err))
	}

	slog.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})

	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	cfg, err := configs.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	echoServer := server.NewServer(*cfg)
	echoServer.SetDefaultMiddlewares()

	err = echoServer.RegisterRoutes()
	if err != nil {
		slog.Error("failed to register routes", slog.Any("error", err))
		panic(err)
	}

	httpServer := echoServer.ToHTTPServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(httpServer, done, cfg.HTTP.ShutdownTimeoutInSeconds)

	slog.Info("starting server", slog.String("address", httpServer.Addr))
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", slog.Any("error", err))
		panic(err)
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("graceful shutdown complete")
}
