package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RunOptions struct {
	Handlers []Handler
}

type Handler struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}

func Init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func Run(options RunOptions) {
	r := gin.Default()

	for _, handler := range options.Handlers {
		if handler.Path == "" {
			handler.Path = "/"
		}

		r.Handle(handler.Method, handler.Path, handler.Handler)
	}

	r.Run("0.0.0.0:80")
}
