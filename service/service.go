package service

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"text/template"

	"github.com/MaxIvanyshen/budget-tracker/database/sqlc"
	"github.com/MaxIvanyshen/budget-tracker/types"
	"github.com/gorilla/sessions"
)

//go:embed templates
var templates embed.FS

type Service struct {
	router       *http.ServeMux
	logger       *slog.Logger
	tmpl         *template.Template
	queries      *sqlc.Queries
	sessionStore *sessions.CookieStore
}

func Start(router *http.ServeMux, logger *slog.Logger, db *sql.DB) {
	svc := &Service{
		router:  router,
		logger:  logger,
		queries: sqlc.New(db),
	}

	svc.sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	tmpl := template.New("base")

	// Add functions once
	tmpl = tmpl.Funcs(template.FuncMap{
		"toJSON": func(v any) string {
			b, err := json.Marshal(v)
			if err != nil {
				// Handle error appropriately
				return ""
			}
			return string(b)
		},
	})

	// Parse templates
	tmpl = template.Must(tmpl.ParseFS(templates, "templates/*.htmx"))

	svc.tmpl = tmpl

	routes := svc.routes()
	for _, route := range routes {
		methodAndPath := route.Method + " " + route.Path
		handler := route.Handler
		if route.Auth {
			handler = authMiddleware(handler, svc.sessionStore).ServeHTTP
		}
		router.HandleFunc(methodAndPath, handler)
		svc.logger.LogAttrs(context.Background(), slog.LevelInfo, "Registered route", slog.String("methodAndPath", methodAndPath), slog.Bool("auth", route.Auth))
	}
}

func (s *Service) runTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	w.Header().Set("Content-Type", "text/html")
	s.tmpl.Funcs(template.FuncMap{
		"toJSON": func(v any) string {
			b, err := json.Marshal(v)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to marshal to JSON", slog.Any("error", err))
				return ""
			}
			return string(b)
		},
	})
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to execute template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func authMiddleware(next types.HandlerFunc, sessions *sessions.CookieStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Get(r, "auth-session")
		if err != nil {
			slog.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.HandlerFunc(next).ServeHTTP(w, r)
	})
}
