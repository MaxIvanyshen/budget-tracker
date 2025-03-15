package service

import (
	"context"
	"embed"
	"log/slog"
	"net/http"
	"text/template"
)

//go:embed templates
var templates embed.FS

type handlerFunc func(http.ResponseWriter, *http.Request)

type Service struct {
	router *http.ServeMux
	logger *slog.Logger
	tmpl   *template.Template
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
		handler := route.Handler
		if route.Auth {
			handler = authMiddleware(handler).ServeHTTP
		}
		router.HandleFunc(methodAndPath, handler)
		svc.logger.LogAttrs(context.Background(), slog.LevelInfo, "Registered route", slog.String("methodAndPath", methodAndPath), slog.Bool("auth", route.Auth))
	}
}

func (s *Service) runTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	w.Header().Set("Content-Type", "text/html")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to execute template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func authMiddleware(next handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		slog.LogAttrs(r.Context(), slog.LevelInfo, "Authorization header", slog.String("header", authHeader))

		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		//TODO: Implement JWT validation

		http.HandlerFunc(next).ServeHTTP(w, r)
	})
}
