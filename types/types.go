package types

import (
	"fmt"
	"net/http"

	"github.com/MaxIvanyshen/budget-tracker/database/sqlc"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

var (
	SupportFormSuccess = `
  <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4" role="alert">
      <p class="font-bold">Success!</p>
      <p>Your message has been sent. We'll get back to you shortly.</p>
  </div>
  <div id="fullname" hx-swap-oob="true"><input type="text" name="fullname" value="" /></div>
  <div id="email" hx-swap-oob="true"><input type="email" name="email" value="" /></div>
  <div id="subject" hx-swap-oob="true"><input type="text" name="subject" value="" /></div>
  <div id="message" hx-swap-oob="true"><textarea name="message"></textarea></div>
`
	Error = `
    <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4" role="alert">
    <p class="font-bold">%s</p>
    </div>
    `
)

func Redirect(path string) string {
	return fmt.Sprintf(`
        <script>
            window.location.href = "%s";
        </script>
        `,
		path,
	)
}

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
	Auth    bool
}

type Data struct {
	Paths          []string
	User           *sqlc.Users
	AdditionalData map[string]any
}

type SupportMsg struct {
	Name    string
	Email   string
	Subject string
	Message string
}
