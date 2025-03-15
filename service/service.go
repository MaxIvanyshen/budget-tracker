package service

import (
	"embed"
	"log/slog"
	"net/http"
	"text/template"
)

//go:embed templates
var templates embed.FS

type Service struct {
	router *http.ServeMux
	logger *slog.Logger
	tmpl   *template.Template
}

type Route struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

func Start(router *http.ServeMux, logger *slog.Logger) {
	svc := &Service{
		router: router,
		logger: logger,
	}

	svc.tmpl = template.Must(template.ParseFS(templates, "templates/*.htmx"))

	routes := svc.routes()
	for _, route := range routes {
		methodAndPath := route.Method + " " + route.Path
		router.HandleFunc(methodAndPath, route.Handler)
	}
}
