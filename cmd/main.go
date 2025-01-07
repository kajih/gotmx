package main

import (
	"embed"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

//go:embed embed
var embeddedFiles embed.FS

func main() {
	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "/embed/web", // because files are located in `web` directory in `embed` fs
		Filesystem: http.FS(embeddedFiles),
	}))

	api := e.Group("/api")
	api.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "users")
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
