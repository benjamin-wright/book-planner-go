package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/tokens/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/internal/tokens"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func GetLoginToken(t tokens.Keyfile) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/:subject/login",
		Handler: func(c *gin.Context) {
			subject := c.Param("subject")

			loginToken, err := t.New(subject, "login", time.Hour)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to create login token: %+v", err))
				return
			}

			c.JSON(http.StatusOK, client.GetLoginTokenResponse{
				Token:  loginToken,
				MaxAge: int(time.Hour / time.Second),
			})
		},
	}
}
