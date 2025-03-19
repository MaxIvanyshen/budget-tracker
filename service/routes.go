package service

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: s.handleUserLogin,
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
			Path:    "/income",
			Handler: s.handleIncome,
			Auth:    true,
		},
		{
			Method:  http.MethodGet,
			Path:    "/expenses",
			Handler: s.handleExpenses,
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

func (s *Service) buildData(w http.ResponseWriter, r *http.Request) types.Data {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Building data", slog.String("path", r.URL.Path))
	data := types.Data{
		Paths:          s.getPaths(),
		AdditionalData: make(map[string]any),
	}
	session, err := s.sessionStore.Get(r, "auth-session")
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get session", slog.Any("error", err))
	}
	id, ok := session.Values["userId"].(int64)
	if !ok {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "No user ID in session")
		if r.URL.Path != "/login" && r.URL.Path != "/signup" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return data
		}
	}
	user, err := s.queries.GetUserByID(r.Context(), id)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to get user", slog.Any("error", err))
		return data
	}

	data.User = user

	return data
}

func (s *Service) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "index", s.buildData(w, r))
}

func (s *Service) handleSignup(w http.ResponseWriter, r *http.Request) {
	data := s.buildData(w, r)
	if data.User != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.runTemplate(w, r, "signup", data)
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
	w.Write([]byte(types.Redirect("/dashboard")))
}

func (s *Service) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := s.buildData(w, r)
	data.AdditionalData["firstname"] = strings.Split(data.User.Name, " ")[0]
	s.runTemplate(w, r, "dashboard", data)
}

func (s *Service) handleContact(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "contact", s.buildData(w, r))
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
	data := s.buildData(w, r)
	if data.User != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	s.runTemplate(w, r, "login", data)
}

func (s *Service) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to parse form", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user, err := s.loginUser(r.Context(), r.Form)
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
	w.Write([]byte(types.Redirect("/dashboard")))
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

func (s *Service) handleIncome(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "income", s.buildData(w, r))
}

func (s *Service) handleExpenses(w http.ResponseWriter, r *http.Request) {
	s.runTemplate(w, r, "expenses", s.buildData(w, r))
}
