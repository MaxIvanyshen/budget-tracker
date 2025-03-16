package types

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
	Auth    bool
}

type Data struct {
	Paths []string
	User  User
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type SupportMsg struct {
	Name    string
	Email   string
	Subject string
	Message string
}
