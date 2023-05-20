package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/apis/books/pkg/client"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/api"
)

func GetBooks(c *database.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/user/:user/books",
		Handler: func(ctx *gin.Context) {
			user := ctx.Param("user")
			zap.S().Infof("Getting books for user %s", user)

			books, err := c.GetBooks(user)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			response := client.GetBooksResponse{
				Books: make([]client.Book, 0, len(books)),
			}

			for _, book := range books {
				response.Books = append(response.Books, client.Book{
					ID:          book.ID,
					Name:        book.Name,
					Summary:     book.Summary,
					CreatedTime: book.CreatedTime,
				})
			}

			ctx.JSON(http.StatusOK, &response)
		},
	}
}
