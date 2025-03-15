package service

import (
	"net/http"
)

func (s *Service) routes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: s.handleIndex,
			Auth:    false,
		},
		{
			Method:  http.MethodGet,
			Path:    "/signup",
			Handler: s.handleSignup,
			Auth:    false,
		},
	}
}

func (s *Service) getPaths() []string {
	routes := s.routes()
	paths := make([]string, 0, len(routes))
	for _, route := range routes {
		if route.Method == http.MethodGet {
			paths = append(paths, route.Path)
		}
	}
	return paths
}

func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "index", Data{Paths: s.getPaths()})
}

func (s *Service) handleSignup(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "signup", Data{Paths: s.getPaths()})
}
