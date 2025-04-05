package server

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterRoutes() error {
	// Register the frontend route
	err := NewFrontendRoute(s.cfg, s.BaseEchoServer)
	if err != nil {
		slog.Error("failed to register frontend route", slog.Any("error", err))
		return err
	}

	// Register the API routes
	s.APIGroup.GET("/hello", s.HelloWorldHandler)

	return nil
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
