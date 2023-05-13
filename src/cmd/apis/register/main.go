package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/register/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	keyfile := os.Getenv("TOKEN_KEYFILE")

	t, err := tokens.New([]byte(keyfile))
	if err != nil {
		zap.S().Errorf("Failed to create token client: %+v", err)
	}

	api.Run(api.RunOptions{
		PostHandler: func(c *gin.Context) {
			var body client.RegisterPostBody

			err := c.BindJSON(&body)
			if err != nil {
				c.AbortWithError(400, fmt.Errorf("failed to parse request body: %+v", err))
				return
			}

			err = t.AddPasswordHash(body.Username, body.Password)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to add password hash: %+v", err))
				return
			}

			c.Status(http.StatusAccepted)
		},
	})
}
