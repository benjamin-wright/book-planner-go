package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/api/request"
)

type Client struct {
	url string
}

func New(URL string) *Client {
	return &Client{
		url: URL,
	}
}

type Book struct {
	ID          string    `json:"id" validate:"uuid"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	CreatedTime time.Time `json:"createdTime"`
}

type GetBooksResponse struct {
	Books []Book `json:"books"`
}

func (c *Client) GetBooks(user string) (*GetBooksResponse, error) {
	var response GetBooksResponse
	status, err := request.Get(context.TODO(), c.url+"/books", &response)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %+v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", status)
	}

	return &response, nil
}
