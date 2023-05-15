package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/auth/login/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	keyfile := tokens.Keyfile(os.Getenv("TOKEN_KEYFILE"))

	t, err := tokens.New()
	if err != nil {
		zap.S().Errorf("Failed to create token client: %+v", err)
	}

	api.Run(api.RunOptions{
		PostHandler: func(c *gin.Context) {
			var body client.PostBody

			err := c.BindJSON(&body)
			if err != nil {
				c.AbortWithError(400, fmt.Errorf("failed to parse request body: %+v", err))
				return
			}

			ok, err := t.CheckPassword(body.Username, body.Password)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to check password: %+v", err))
				return
			}

			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			loginToken, err := keyfile.New(body.Username, "login", time.Hour)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to create login token: %+v", err))
				return
			}

			c.JSON(http.StatusCreated, client.TokenResponse{
				Token:  loginToken,
				MaxAge: int(time.Hour / time.Second),
			})
		},
	})
}
