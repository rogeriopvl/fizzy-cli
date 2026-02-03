package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	endpointURL := c.AccountBaseURL + "/users"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response []User
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	endpointURL := c.AccountBaseURL + "/users/" + userID

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response User
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PutUser(ctx context.Context, userID string, payload UpdateUserPayload) error {
	endpointURL := c.AccountBaseURL + "/users/" + userID

	body := map[string]UpdateUserPayload{"user": payload}

	req, err := c.newRequest(ctx, http.MethodPut, endpointURL, body)
	if err != nil {
		return fmt.Errorf("failed to create update user request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	endpointURL := c.AccountBaseURL + "/users/" + userID

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}
