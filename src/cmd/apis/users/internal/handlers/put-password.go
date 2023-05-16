package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func PutPassword(c *database.Client) api.Handler {
	return api.Handler{
		Method: "PUT",
		Path:   "/:name/password",
		Handler: func(ctx *gin.Context) {
			name := ctx.Param("name")
			var body client.CheckPasswordRequest
			err := ctx.Bind(&body)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			user, err := c.CheckPassword(name, body.Password)
			if err == database.ErrPasswordMismatch {
				ctx.AbortWithError(http.StatusUnauthorized, err)
				return
			} else if err != nil {
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
