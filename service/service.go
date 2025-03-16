package service

import (
	"context"
	"database/sql"
	"embed"
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

	svc.tmpl = template.Must(template.ParseFS(templates, "templates/*.htmx"))

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

func (s *Service) getUserID(w http.ResponseWriter, r *http.Request) (int64, error) {
	session, err := s.sessionStore.Get(r, "auth-session")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	userID, ok := session.Values["userId"].(int64)
	if !ok {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get user_id", slog.Any("error", err))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	return userID, nil

}
