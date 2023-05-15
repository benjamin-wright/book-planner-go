package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	hostname := os.Getenv("WEB_HOSTNAME")
	loginURL := os.Getenv("LOGIN_URL")
	proxyPrefix := os.Getenv("PROXY_PREFIX")

	api.Run(api.RunOptions{
		Path: proxyPrefix,
		GetHandler: func(c *gin.Context) {
			c.SetCookie("ponglehub.login", "", -1, "/", hostname, false, true)
			c.Redirect(http.StatusTemporaryRedirect, loginURL)
		},
	})
}
