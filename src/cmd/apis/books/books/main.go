package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func main() {
	api.Init()

	_, err := client.New()
	if err != nil {
		zap.S().Fatalf("Failed to create postgres client: %+v", err)
	}

	api.Run(api.RunOptions{
		Handlers: []api.Handler{
			{
				Method: "GET",
				Handler: func(c *gin.Context) {
					c.Status(200)
				},
			},
		},
	})
}
