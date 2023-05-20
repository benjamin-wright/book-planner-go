package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	api.Init()

	hostname := os.Getenv("WEB_HOSTNAME")
	loginURL := os.Getenv("LOGIN_URL")
	proxyPrefix := os.Getenv("PROXY_PREFIX")

	router := api.Router(api.RunOptions{
		Handlers: []api.Handler{
			{
				Path:   proxyPrefix,
				Method: "GET",
				Handler: func(c *gin.Context) {
					c.SetCookie("ponglehub.login", "", -1, "/", hostname, false, true)
					c.Redirect(http.StatusTemporaryRedirect, loginURL)
				},
			},
		},
	})

	api.Run(router)
}
