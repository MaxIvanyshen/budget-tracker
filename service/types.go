package service

type Route struct {
	Method  string
	Path    string
	Handler handlerFunc
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
