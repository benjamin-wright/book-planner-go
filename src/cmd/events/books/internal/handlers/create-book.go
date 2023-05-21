package handlers

import (
	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/books/internal/database"
	"ponglehub.co.uk/book-planner-go/src/cmd/events/pkg/types"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/events"
)

func CreateBook(client *database.Client, subject string) events.RequestHandler[types.CreateBookEvent, types.CreateBookResponseEvent] {
	return events.RequestHandler[types.CreateBookEvent, types.CreateBookResponseEvent]{
		Channel: subject,
		Handler: func(event *types.CreateBookEvent) types.CreateBookResponseEvent {
			zap.S().Infof("[%s] Creating book", event.UserID)

			response := types.CreateBookResponseEvent{
				UserID: event.UserID,
				Name:   event.Name,
			}

			err := client.AddBook(database.Book{
				UserID:  event.UserID,
				Name:    event.Name,
				Summary: event.Summary,
			})

			if err == database.ErrAlreadyExists {
				zap.S().Warnf("[%s] Already exists", event.UserID)
				response.Status = types.AlreadyExists
			} else if err != nil {
				zap.S().Errorf("[%s] Server error: %+v", event.UserID, err)
				response.Status = types.ServerError
			} else {
				zap.S().Infof("[%s] Created", event.UserID)
				response.Status = types.Created
			}

			return response
		},
	}
}
