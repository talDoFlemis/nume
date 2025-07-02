package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/taldoflemis/nume/configs"
	"github.com/taldoflemis/nume/internal/tui/models"
)

func gracefulShutdown(
	s *ssh.Server,
	done chan bool,
	shutdownTimeoutInSeconds int,
) {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	slog.Info("shutting down gracefully. press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(shutdownTimeoutInSeconds)*time.Second,
	)
	defer cancel()
	slog.Info("server exiting")

	// Shutdown the server gracefully
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server gracefully", slog.Any("error", err))
		return
	}

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

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(cfg.SSH.Host, strconv.Itoa(cfg.SSH.Port))),
		wish.WithHostKeyPath(cfg.SSH.HostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.StructuredMiddleware(),
		),
	)
	if err != nil {
		slog.Error("failed to create SSH server", slog.Any("error", err))
		return
	}

	done := make(chan bool)
	go gracefulShutdown(s, done, cfg.HTTP.ShutdownTimeoutInSeconds)

	slog.Info("starting SSH server")

	err = s.ListenAndServe()
	if err != nil {
		slog.Error("failed to start SSH server", slog.Any("error", err))
		return
	}

	slog.Info("SSH server started")

	// Wait for the shutdown signal
	<-done
	slog.Info("SSH server down")
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	opts := bubbletea.MakeOptions(s)
	opts = append(opts, tea.WithAltScreen())

	theme := models.ThemeCatppuccin(renderer)
	m := models.NewWelcomeModel(theme, pty.Term, renderer.ColorProfile().Name(), s.User())
	return m, opts
}
