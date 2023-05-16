package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"ponglehub.co.uk/book-planner-go/src/pkg/web/api/request"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{url: url}
}

type GetLoginTokenResponse struct {
	Token  string `json:"token"`
	MaxAge int    `json:"maxAge"`
}

func (c *Client) GetLoginToken(subject string) (*GetLoginTokenResponse, error) {
	var response GetLoginTokenResponse
	status, err := request.Get(context.TODO(), c.url+"/"+subject+"/login", &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", status)
	}

	return &response, nil
}

type ValidateLoginTokenResponse struct {
	Subject string `json:"subject"`
}

func (c *Client) ValidateLoginToken(token string) (*ValidateLoginTokenResponse, error) {
	var response ValidateLoginTokenResponse
	status, err := request.Get(context.TODO(), c.url+"/validate/login?token="+url.PathEscape(token), &response)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", status)
	}

	return &response, nil
}
