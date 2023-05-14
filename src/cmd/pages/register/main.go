package main

import (
	_ "embed"
	"errors"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/auth/register/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api/validation"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

type Context struct {
	SubmitURL string
}

func main() {
	baseURL := os.Getenv("BASE_URL")
	proxyPrefix := os.Getenv("PROXY_PREFIX")
	submitURL := os.Getenv("SUBMIT_URL")
	redirectURL := os.Getenv("REDIRECT_URL")

	cli := client.New(os.Getenv("REGISTER_API_URL"))

	runtime.Run(runtime.ServerOptions{
		Template: content,
		Title:    "Book Planner",
		PageHandler: func(r *http.Request) any {
			return Context{
				SubmitURL: submitURL,
			}
		},
		PostHandler: func(w http.ResponseWriter, r *http.Request) {
			zap.S().Info("register submission")

			err := r.ParseForm()
			if err != nil {
				zap.S().Errorf("error parsing form: %+v", err)
				http.Redirect(w, r, baseURL+proxyPrefix+"?error=unknown", http.StatusFound)
				return
			}

			missing := validation.GetMissingFields(r.Form, []string{"username", "password", "confirm-password"})
			if len(missing) > 0 {
				zap.S().Warnf("missing fields: %+v", missing)
				err := make([]string, 0, len(missing))
				for _, field := range missing {
					err = append(err, "missing="+field)
				}
				http.Redirect(w, r, baseURL+proxyPrefix+"?"+strings.Join(err, "&"), http.StatusFound)
				return
			}

			password := r.Form.Get("password")
			confirm := r.Form.Get("confirm-password")

			if password != confirm {
				zap.S().Warn("mistmatched passwords")
				http.Redirect(w, r, baseURL+proxyPrefix+"?error=password", http.StatusFound)
				return
			}

			if !validation.CheckPasswordComplexity(password) {
				zap.S().Warn("password complexity")
				http.Redirect(w, r, baseURL+proxyPrefix+"?error=complexity", http.StatusFound)
				return
			}

			username := r.Form.Get("username")

			zap.S().Infof("Adding new user %s with password %s", username, password)

			err = cli.Register(r.Context(), client.PostBody{
				Username: username,
				Password: password,
			})
			if err != nil {
				zap.S().Errorf("error sending registration request: %+v", err)

				if errors.Is(err, client.UserExistsError) {
					http.Redirect(w, r, baseURL+proxyPrefix+"?error=exists", http.StatusFound)
				} else {
					http.Redirect(w, r, baseURL+proxyPrefix+"?error=unknown", http.StatusFound)
				}

				return
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
		},
	})
}
