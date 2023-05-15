package main

import (
	_ "embed"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/auth/login/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api/validation"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

//go:embed module.wasm
var wasm []byte

type Context struct {
	RegisterURL string
	Registered  bool
	SubmitURL   string
}

func main() {
	registerURL := os.Getenv("REGISTER_URL")
	hostname := os.Getenv("WEB_HOSTNAME")
	proxyPrefix := os.Getenv("PROXY_PREFIX")
	submitURL := os.Getenv("SUBMIT_URL")
	redirectURL := os.Getenv("REDIRECT_URL")

	cli := client.New(os.Getenv("LOGIN_API_URL"))

	runtime.Run(runtime.ServerOptions{
		Template:    content,
		Title:       "Book Planner: Login",
		HideHeaders: true,
		WASMModules: []runtime.WASMModule{
			{
				Name: "wasm",
				Path: "module.wasm",
				Data: wasm,
			},
		},
		PageHandler: func(r *http.Request) any {
			query := r.URL.Query()
			registered := query.Has("registered")

			return Context{
				RegisterURL: registerURL,
				Registered:  registered,
				SubmitURL:   submitURL,
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

			missing := validation.GetMissingFields(r.Form, []string{"username", "password"})
			if len(missing) > 0 {
				zap.S().Warnf("missing fields: %+v", missing)
				err := make([]string, 0, len(missing))
				for _, field := range missing {
					err = append(err, "missing="+field)
				}
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?"+strings.Join(err, "&"), http.StatusFound)
				return
			}

			username := r.Form.Get("username")
			password := r.Form.Get("password")

			zap.S().Infof("Logging in user %s", username)

			response, err := cli.Login(r.Context(), client.PostBody{
				Username: username,
				Password: password,
			})
			if err != nil {
				zap.S().Errorf("error sending login request: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unauthorized", http.StatusFound)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "ponglehub.login",
				Value:    response.Token,
				Domain:   hostname,
				MaxAge:   response.MaxAge,
				HttpOnly: true,
			})

			http.Redirect(w, r, redirectURL, http.StatusFound)
		},
	})
}
