package service

import (
	"log/slog"
	"net/http"

	"github.com/MaxIvanyshen/budget-tracker/types"
	"github.com/MaxIvanyshen/budget-tracker/utils"
)

func (s *Service) routes() []types.Route {
	return []types.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: s.handleIndex,
		},
		{
			Method:  http.MethodGet,
			Path:    "/signup",
			Handler: s.handleSignup,
		},
		{
			Method:  http.MethodPost,
			Path:    "/signup",
			Handler: s.handleUserRegistration,
		},
		{
			Method:  http.MethodGet,
			Path:    "/dashboard",
			Handler: s.handleDashboard,
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/contact",
			Handler: s.handleContact,
		},
		{
			Method:  http.MethodPost,
			Path:    "/send-contact-form",
			Handler: s.handleSendContactMsg,
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
	s.runTemplate(w, r, "index", types.Data{Paths: s.getPaths()})
}

func (s *Service) handleSignup(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "signup", types.Data{Paths: s.getPaths()})
}

func (s *Service) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement user registration and add access token to response
	r.Header.Set("Authorization", "Bearer access-token")
	s.handleDashboard(w, r)
}

func (s *Service) handleDashboard(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "dashboard", types.Data{Paths: s.getPaths()})
}

func (s *Service) handleContact(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "contact", types.Data{Paths: s.getPaths()})
}

func (s *Service) handleSendContactMsg(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Processing contact form submission")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to parse form", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	msg := types.SupportMsg{
		Name:    r.FormValue("fullname"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	err = utils.SendTelegramMessage(msg)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to send message", slog.Any("error", err))
	}

	successHTML := `
  <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4" role="alert">
      <p class="font-bold">Success!</p>
      <p>Your message has been sent. We'll get back to you shortly.</p>
  </div>
  <div id="fullname" hx-swap-oob="true"><input type="text" name="fullname" value="" /></div>
  <div id="email" hx-swap-oob="true"><input type="email" name="email" value="" /></div>
  <div id="subject" hx-swap-oob="true"><input type="text" name="subject" value="" /></div>
  <div id="message" hx-swap-oob="true"><textarea name="message"></textarea></div>
`
	w.Write([]byte(successHTML))
}
