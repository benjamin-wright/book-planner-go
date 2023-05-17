package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Get(ctx context.Context, url string, response any) (int, error) {
	return request(ctx, "GET", url, nil, response)
}

func Put(ctx context.Context, url string, body any, response any) (int, error) {
	return request(ctx, "PUT", url, body, response)
}

func Post(ctx context.Context, url string, body any, response any) (int, error) {
	return request(ctx, "POST", url, body, response)
}

func request(ctx context.Context, method string, url string, body any, response any) (int, error) {
	postData := &bytes.Buffer{}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal post body: %+v", err)
		}

		postData = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, postData)
	if err != nil {
		return 0, fmt.Errorf("failed to create request object: %+v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %+v", err)
	}
	defer res.Body.Close()

	if response != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&response)
		if err != nil {
			return res.StatusCode, fmt.Errorf("failed to decode response: %+v", err)
		}
	}

	return res.StatusCode, nil
}
