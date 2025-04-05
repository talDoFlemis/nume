package server

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/taldoflemis/nume/configs"
	"github.com/taldoflemis/nume/frontend"
)

func NewFrontendRoute(cfg configs.Config, e *echo.Echo) error {
	if cfg.App.Environment == "local" {
		return setupViteDevProxy(cfg.HTTP, e)
	}

	setupViteProd(cfg.HTTP, e)

	return nil
}

func setupViteProd(cfg configs.HTTPCfg, e *echo.Echo) {
	// Setup MPA serving
	// Use the static assets from the dist directory
	e.FileFS("/", "index.html", frontend.DistIndexHTML)
	e.StaticFS("/", frontend.DistDirFS)

	// This is needed to serve the index.html file for all routes that
	// are not /api/* needed for SPA to work when loading a specific
	// url directly
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			// Skip the proxy if the prefix is /api
			return len(c.Path()) >= len(cfg.APIPrefix) &&
				c.Path()[:len(cfg.APIPrefix)] == cfg.APIPrefix
		},
		// Root directory from where the static content is served.
		Root: "/",
		// Enable HTML5 mode by forwarding all not-found requests to
		// root so that SPA (single-page application) can handle the routing.
		HTML5:      false,
		Browse:     false,
		IgnoreBase: true,
		Filesystem: http.FS(frontend.DistDirFS),
	}))
}

func setupViteDevProxy(cfg configs.HTTPCfg, e *echo.Echo) error {
	urlParsed, err := url.Parse("http://localhost:5173")
	if err != nil {
		slog.Error("failed to parse url for dev proxy", slog.Any("error", err))
		return err
	}

	// Setep a proxy to the vite dev server on localhost:5173
	balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: urlParsed,
		},
	})

	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: balancer,
		Skipper: func(c echo.Context) bool {
			// Skip the proxy if the prefix is /api
			return len(c.Path()) >= 4 && c.Path()[:4] == cfg.APIPrefix
		},
	}))

	return nil
}
