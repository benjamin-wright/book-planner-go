package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PostBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token  string `json:"token"`
	MaxAge int    `json:"maxAge"`
}

type Client struct {
	url string
}

func New(URL string) *Client {
	return &Client{
		url: URL,
	}
}

func (c *Client) Login(ctx context.Context, body PostBody) (*TokenResponse, error) {
	postData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal post body: %+v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request object: %+v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %+v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var response TokenResponse
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %+v", err)
	}

	return &response, nil
}
