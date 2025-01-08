package main

import (
	"embed"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
)

//go:embed embed
var embeddedFiles embed.FS

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseFS(embeddedFiles, "embed/templates/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Count struct {
	Count int
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "/embed/web", // because files are located in `web` directory in `embed` fs
		Filesystem: http.FS(embeddedFiles),
	}))

	count := Count{Count: 0}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "blocks-index", count)
	})

	e.GET("/getcount", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "countBlock", count)
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
