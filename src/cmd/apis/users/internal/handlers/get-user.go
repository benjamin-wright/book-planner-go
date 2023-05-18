package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func GetUser(c *database.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/:name",
		Handler: func(ctx *gin.Context) {
			name := ctx.Param("name")
			user, err := c.GetUser(name)
			if err == database.ErrNoUser {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			} else if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusOK, client.GetUserResponse{
				ID:       user.ID,
				Username: user.Name,
			})
		},
	}
}
