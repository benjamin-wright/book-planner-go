package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/tokens"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	t, err := tokens.New([]byte{})
	if err != nil {
		zap.S().Errorf("Failed to create token client: %+v", err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		token, err := c.Cookie("ponglehub.login")
		if err != nil {
			zap.S().Info("No cookie token")
			c.Redirect(http.StatusTemporaryRedirect, loginURL)
			return
		}

		claims, err := t.Parse(token)
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
	})

	r.Run("0.0.0.0:80")
}
