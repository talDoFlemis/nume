package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/taldoflemis/nume/configs"
)

type Server struct {
	port           int
	BaseEchoServer *echo.Echo
	cfg            configs.Config
	APIGroup       *echo.Group
}

func NewServer(httpConfig configs.Config) *Server {
	e := echo.New()
	api := e.Group(httpConfig.HTTP.ApiPrefix)

	NewServer := &Server{
		port:           httpConfig.HTTP.Port,
		BaseEchoServer: e,
		cfg:            httpConfig,
		APIGroup:       api,
	}

	return NewServer
}

func (s *Server) SetDefaultMiddlewares() {
	s.BaseEchoServer.IPExtractor = echo.ExtractIPFromXFFHeader()
	s.BaseEchoServer.Use(slogecho.New(slog.Default()))
	s.BaseEchoServer.Use(middleware.Recover())
	s.BaseEchoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     s.cfg.HTTP.CORS.Origins,
		AllowMethods:     s.cfg.HTTP.CORS.Methods,
		AllowHeaders:     s.cfg.HTTP.CORS.Headers,
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (s *Server) ToHttpServer() *http.Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.BaseEchoServer,
		IdleTimeout:  time.Duration(s.cfg.HTTP.IdleTimeoutInSeconds) * time.Second,
		ReadTimeout:  time.Duration(s.cfg.HTTP.ReadTimeoutInSeconds) * time.Second,
		WriteTimeout: time.Duration(s.cfg.HTTP.WriteTimeoutInSeconds) * time.Second,
	}

	return server
}
