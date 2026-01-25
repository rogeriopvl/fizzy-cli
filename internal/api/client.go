// Package api provides the Fizzy API client
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL        string
	AccountBaseURL string
	BoardBaseURL   string
	AccessToken    string
	HTTPClient     *http.Client
}

func NewClient(accountSlug string, boardID string) (*Client, error) {
	baseURL := "https://app.fizzy.do"
	accountBaseURL := baseURL + accountSlug

	var boardBaseURL string
	if boardID != "" {
		boardBaseURL = accountBaseURL + "/boards" + "/" + boardID
	}

	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return nil, fmt.Errorf("FIZZY_ACCESS_TOKEN environment variable is not set")
	}

	return &Client{
		BaseURL:        baseURL,
		AccountBaseURL: accountBaseURL,
		BoardBaseURL:   boardBaseURL,
		AccessToken:    token,
		HTTPClient:     &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// newRequest makes an HTTP request with the required headers setup
func (c *Client) newRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// decodeResponse executes a request and decodes the JSON response into v
// If expectedStatus is 0, it defaults to http.StatusOK
// If v is nil, the response body is not decoded
func (c *Client) decodeResponse(req *http.Request, v any, expectedStatus ...int) (int, error) {
	expectedCode := http.StatusOK
	if len(expectedStatus) > 0 {
		expectedCode = expectedStatus[0]
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != expectedCode {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, fmt.Errorf("unexpected status code %d (failed to read error response: %w)", res.StatusCode, err)
		}
		return 0, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(res.Body).Decode(v); err != nil {
			return 0, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return res.StatusCode, nil
}
