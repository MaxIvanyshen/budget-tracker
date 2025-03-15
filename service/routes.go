package service

import (
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Service) routes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: s.handleIndex,
		},
		{
			Method:  http.MethodPost,
			Path:    "/counter",
			Handler: s.handleCounter,
		},
	}
}

func (s *Service) runTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	w.Header().Set("Content-Type", "text/html")
	if err := s.tmpl.ExecuteTemplate(w, name, data); err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to execute template", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type data struct {
	Name    string `json:"name"`
	Counter int    `json:"counter"`
}

func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "index", data{Name: "World", Counter: 0})
}

func (s *Service) handleCounter(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	count, _ := strconv.Atoi(r.FormValue("count"))
	count++
	s.runTemplate(w, r, "counter", data{Name: "World", Counter: count})
}
