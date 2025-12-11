// Package api
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	baseURL     string
	accessToken string
}

func NewClient(accountSlug string) (*Client, error) {
	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return nil, fmt.Errorf("FIZZY_ACCESS_TOKEN environment variable is not set")
	}

	return &Client{
		baseURL:     fmt.Sprintf("https://app.fizzy.do%s", accountSlug),
		accessToken: token,
	}, nil
}

// newRequest makes an HTTP request with the required headers setup
func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// decodeResponse executes a request and decodes the JSON response into v
func (c *Client) decodeResponse(req *http.Request, v any) error {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}

	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func (c *Client) GetBoards(ctx context.Context) ([]Board, error) {
	endpointURL := c.baseURL + "/boards"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response []Board
	if err := c.decodeResponse(req, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetMyIdentity(ctx context.Context) (*GetMyIdentityResponse, error) {
	endpointURL := c.baseURL + "/my/identity"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response GetMyIdentityResponse
	if err := c.decodeResponse(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type Board struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AllAccess bool   `json:"all_access"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
	Creator   User   `json:"creator"`
}

type GetMyIdentityResponse struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      User   `json:"user"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email_address"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
}
