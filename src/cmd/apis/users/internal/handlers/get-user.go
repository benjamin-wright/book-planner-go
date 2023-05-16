package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func GetUser(c *database.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/:name",
		Handler: func(ctx *gin.Context) {
			name := ctx.Param("name")
			user, err := c.GetUser(name)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusOK, map[string]string{
				"username": user.Name,
				"id":       user.ID,
			})
		},
	}
}
