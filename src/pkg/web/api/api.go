package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RunOptions struct {
	GetHandler  func(c *gin.Context)
	PostHandler func(c *gin.Context)
}

func Run(options RunOptions) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	r := gin.Default()
	if options.GetHandler != nil {
		r.GET("/", options.GetHandler)
	}

	if options.PostHandler != nil {
		r.POST("/", options.PostHandler)
	}

	r.Run("0.0.0.0:80")
}
