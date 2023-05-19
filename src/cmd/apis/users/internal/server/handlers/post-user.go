package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func PostUser(c *database.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   "/",
		Handler: func(ctx *gin.Context) {
			var body client.AddUserRequest
			err := ctx.BindJSON(&body)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			err = c.AddUser(database.User{Name: body.Username, Password: body.Password})
			if err == database.ErrUserExists {
				ctx.AbortWithError(http.StatusConflict, err)
				return
			} else if err == database.ErrComplexity {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			} else if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.Status(http.StatusCreated)
		},
	}
}
