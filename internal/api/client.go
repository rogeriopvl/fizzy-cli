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
	clientSlug  string
}

func NewClient() (*Client, error) {
	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return nil, fmt.Errorf("FIZZY_ACCESS_TOKEN environment variable is not set")
	}

	return &Client{
		baseURL:     "https://app.fizzy.do",
		accessToken: token,
	}, nil
}

func (c *Client) GetBoards(ctx context.Context) ([]Board, error) {
	return nil, nil
}

func (c *Client) GetMyIdentity(ctx context.Context) (*GetMyIdentityResponse, error) {
	endpointURL := c.baseURL + c.clientSlug + "/my/identity"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}

	var response GetMyIdentityResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

type Board struct {
	ID   string
	Name string
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
