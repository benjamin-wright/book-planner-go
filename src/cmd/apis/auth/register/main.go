package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/auth/register/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/wasm/validation"
)

func main() {
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

			if body.Username == "" {
				c.AbortWithError(400, fmt.Errorf("missing username in request: %+v", err))
				return
			}

			if !validation.CheckPasswordComplexity(body.Password) {
				c.AbortWithError(400, fmt.Errorf("password failed complexity checks: %+v", err))
				return
			}

			existing, err := t.GetToken(body.Username, "password")
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed checking for existing user: %+v", err))
				return
			}
			if existing != "" {
				c.Status(http.StatusConflict)
				return
			}

			err = t.AddPasswordHash(body.Username, body.Password)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to add password hash: %+v", err))
				return
			}

			c.Status(http.StatusCreated)
		},
	})
}
