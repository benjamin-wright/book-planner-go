package main

import (
	_ "embed"
	"net/http"
	"os"

	"go.uber.org/zap"
	tokensApi "ponglehub.co.uk/book-planner-go/src/cmd/apis/tokens/pkg/client"
	usersApi "ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/alert"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/runtime"
)

//go:embed index.html
var content string

type Context struct {
	RegisterURL string
	Registered  bool
	SubmitURL   string
	Error       string
}

func main() {
	registerURL := os.Getenv("REGISTER_URL")
	hostname := os.Getenv("WEB_HOSTNAME")
	proxyPrefix := os.Getenv("PROXY_PREFIX")
	submitURL := os.Getenv("SUBMIT_URL")
	redirectURL := os.Getenv("REDIRECT_URL")

	users := usersApi.New(os.Getenv("USERS_API_URL"))
	tokens := tokensApi.New(os.Getenv("TOKENS_API_URL"))

	runtime.Run(runtime.ServerOptions{
		Template:    content,
		Title:       "Book Planner: Login",
		HideHeaders: true,
		Children: []component.Component{
			alert.Get(),
		},
		PageHandler: func(r *http.Request) any {
			query := r.URL.Query()
			registered := query.Has("registered")

			return Context{
				RegisterURL: registerURL,
				Registered:  registered,
				SubmitURL:   submitURL,
				Error: alert.Lookup(
					query.Get("error"),
					map[string]string{
						"unauthorized": "Login failed, did you use the right username and password?",
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

			username := r.Form.Get("username")
			password := r.Form.Get("password")

			zap.S().Infof("Logging in user %s", username)

			user, valid, err := users.CheckPassword(r.Context(), username, password)
			if err != nil {
				zap.S().Errorf("error checking password: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unknown", http.StatusFound)
				return
			}

			if !valid {
				zap.S().Warnf("password didn't match: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unauthorized", http.StatusFound)
				return
			}

			res, err := tokens.GetLoginToken(user.ID)
			if err != nil {
				zap.S().Errorf("error fetching login token: %+v", err)
				http.Redirect(w, r, "http://"+hostname+proxyPrefix+"?error=unknown", http.StatusFound)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "ponglehub.login",
				Value:    res.Token,
				Domain:   hostname,
				MaxAge:   res.MaxAge,
				HttpOnly: true,
			})

			http.Redirect(w, r, redirectURL, http.StatusFound)
		},
	})
}
