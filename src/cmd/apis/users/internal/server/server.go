package server

import (
	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/users/internal/server/handlers"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func Router(cli *database.Client) *gin.Engine {
	return api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.PostUser(cli),
			handlers.GetUser(cli),
			handlers.PutPassword(cli),
		},
	})
}
