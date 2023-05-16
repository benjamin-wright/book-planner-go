package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	tokensApi "ponglehub.co.uk/book-planner-go/src/cmd/apis/tokens/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func Verify(tokens *tokensApi.Client, loginURL string) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/",
		Handler: func(c *gin.Context) {
			token, err := c.Cookie("ponglehub.login")
			if err != nil {
				zap.S().Info("No cookie token")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			res, err := tokens.ValidateLoginToken(token)
			if err != nil {
				zap.S().Info("token validation failed: %+v", err)
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			c.Header("X-Auth-User", res.Subject)
			c.Status(http.StatusOK)
		},
	}
}
