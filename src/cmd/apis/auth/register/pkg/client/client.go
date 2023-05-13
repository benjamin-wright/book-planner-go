package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PostBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	url string
}

var UserExistsError = errors.New("user already exists")

func New(URL string) *Client {
	return &Client{
		url: URL,
	}
}

func (c *Client) Register(ctx context.Context, body PostBody) error {
	postData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal post body: %+v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(postData))
	if err != nil {
		return fmt.Errorf("failed to create request object: %+v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %+v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		return UserExistsError
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	return nil
}
