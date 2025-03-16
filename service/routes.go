package service

import (
	"fmt"
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
			Path:    "/login",
			Handler: s.handleLogin,
		},
		{
			Method:  http.MethodGet,
			Path:    "/logout",
			Handler: s.handleLogout,
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

func (s *Service) buildData(r *http.Request) types.Data {
	data := types.Data{
		Paths: s.getPaths(),
	}
	session, err := s.sessionStore.Get(r, "auth-session")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
	}
	if id, ok := session.Values["userId"].(int64); ok {
		data.UserID = id
	}
	return data
}

func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "index", s.buildData(r))
}

func (s *Service) handleSignup(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "signup", s.buildData(r))
}

func (s *Service) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to parse form", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user, err := s.registerUser(r.Context(), r.Form)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(types.Error, err.Error())))
		return
	}
	session, err := s.sessionStore.Get(r, "auth-session")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = true
	session.Values["userId"] = user.ID
	session.Save(r, w)
	w.Write([]byte(`
        <script>
            window.location.href = "/dashboard";
        </script>
    `))
}

func (s *Service) handleDashboard(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUserID(w, r)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "User ID", slog.Int64("userID", userID))

	s.runTemplate(w, r, "dashboard", s.buildData(r))
}

func (s *Service) handleContact(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "contact", s.buildData(r))
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

	w.Write([]byte(types.SupportFormSuccess))
}

func (s *Service) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "signup", s.buildData(r))
}

func (s *Service) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, "auth-session")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	session.Values["authenticated"] = false
	delete(session.Values, "userId")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
