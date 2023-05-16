package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/api/request"
)

var ErrUserExists = errors.New("user already exists")

type Client struct {
	url string
}

func New(URL string) *Client {
	return &Client{
		url: URL,
	}
}

type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Client) AddUser(ctx context.Context, username string, password string) error {
	status, err := request.Post(ctx, c.url, AddUserRequest{
		Username: username,
		Password: password,
	}, nil)
	if err != nil {
		return err
	}

	if status == http.StatusConflict {
		return ErrUserExists
	}

	if status != http.StatusCreated {
		return fmt.Errorf("failed with status code %d", status)
	}

	return nil
}

type GetUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (c *Client) GetUser(ctx context.Context, username string) (*GetUserResponse, error) {
	var response GetUserResponse
	status, err := request.Get(ctx, c.url+"/"+url.PathEscape(username), &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusCreated {
		return nil, fmt.Errorf("failed with status code %d", status)
	}

	return &response, nil
}

type CheckPasswordRequest struct {
	Password string `json:"password"`
}

type CheckPasswordResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (c *Client) CheckPassword(ctx context.Context, username string, password string) (*CheckPasswordResponse, bool, error) {
	var response CheckPasswordResponse
	status, err := request.Put(
		ctx,
		fmt.Sprintf("%s/%s/password", c.url, url.PathEscape(username)),
		CheckPasswordRequest{
			Password: password,
		},
		&response,
	)
	if err != nil {
		return nil, false, err
	}

	if status == http.StatusUnauthorized {
		return nil, false, nil
	}

	if status != http.StatusCreated {
		return nil, false, fmt.Errorf("failed with status %d", status)
	}

	return &response, true, nil
}
