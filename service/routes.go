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
		{
			Method:  http.MethodPost,
			Path:    "/signup",
			Handler: s.handleUserRegistration,
			Auth:    false,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dashboard",
			Handler: s.handleDashboard,
			Auth:    true,
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

func (s *Service) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement user registration and add access token to response
	r.Header.Set("Authorization", "Bearer access-token")
	s.handleDashboard(w, r)
}

func (s *Service) handleDashboard(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "dashboard", Data{Paths: s.getPaths()})
}
