package handlers

import (
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/events"
)

type Config struct {
	CreateBooksSubject string
}

func Get(client *database.Client, config Config) []events.HandlerFunc {
	return []events.HandlerFunc{
		CreateBook(client, config.CreateBooksSubject).GetHandler(),
	}
}
