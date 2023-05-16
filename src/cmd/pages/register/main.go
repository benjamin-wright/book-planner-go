package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"

	"go.uber.org/zap"
	usersApi "ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/alert"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

//go:embed module.wasm
var wasm []byte

type Context struct {
	SubmitURL string
	LoginURL  string
	Error     string
}

func main() {
	hostname := os.Getenv("WEB_HOSTNAME")
	proxyPrefix := os.Getenv("PROXY_PREFIX")
	loginURL := os.Getenv("LOGIN_URL")
	submitURL := os.Getenv("SUBMIT_URL")
	redirectURL := os.Getenv("REDIRECT_URL")

	users := usersApi.New(os.Getenv("USERS_API_URL"))

	runtime.Run(runtime.ServerOptions{
		Template:    content,
		Title:       "Book Planner: Register",
		HideHeaders: true,
		WASMModules: []runtime.WASMModule{
			{
				Name: "wasm",
				Path: "module.wasm",
				Data: wasm,
			},
		},
		Children: []component.Component{
			alert.Get(),
		},
		PageHandler: func(r *http.Request) any {
			query := r.URL.Query()

			return Context{
				SubmitURL: submitURL,
				LoginURL:  loginURL,
				Error: alert.Lookup(
					query.Get("error"),
					map[string]string{
						"exists":   "That users already exists, please try another username",
						"password": "Looks like your passwords didn't match, try that again",
					},
				),
			}
		},
		PostHandler: func(w http.ResponseWriter, r *http.Request) {
			zap.S().Info("register submission")

			err := r.ParseForm()
			if err != nil {
				zap.S().Errorf("error parsing form: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unknown", http.StatusFound)
				return
			}

			password := r.Form.Get("password")
			confirm := r.Form.Get("confirm-password")

			if password != confirm {
				zap.S().Warn("mistmatched passwords")
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=password", http.StatusFound)
				return
			}

			username := r.Form.Get("username")

			zap.S().Infof("Adding new user %s", username)
			err = users.AddUser(context.TODO(), username, password)
			if err == usersApi.ErrUserExists {
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=exists", http.StatusFound)
				return
			} else if err != nil {
				zap.S().Errorf("error creating user: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unknown", http.StatusFound)
				return
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
		},
	})
}
