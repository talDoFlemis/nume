package frontend

import (
	"embed"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed dist/*
	dist embed.FS

	//go:embed dist/index.html
	indexHTML embed.FS

	DistDirFS     = echo.MustSubFS(dist, "dist")
	DistIndexHTML = echo.MustSubFS(indexHTML, "dist")
)

