package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")
	keyfile := tokens.Keyfile(os.Getenv("TOKEN_KEYFILE"))

	api.Run(api.RunOptions{
		GetHandler: func(c *gin.Context) {
			token, err := c.Cookie("ponglehub.login")
			if err != nil {
				zap.S().Info("No cookie token")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			claims, err := keyfile.Parse(token)
			if err != nil {
				zap.S().Infof("Failed to parse claims: %+v", err)
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			if claims.Kind != "Login" {
				zap.S().Infof("Non-login claim: %+v", claims)
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			c.Header("X-auth-user", claims.Subject)
		},
	})
}
