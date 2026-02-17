package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetCardComments(ctx context.Context, cardNumber int, opts *ListOptions) ([]Comment, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get card comments request: %w", err)
	}

	limit := 0
	if opts != nil {
		limit = opts.Limit
	}

	return fetchAllPages[Comment](ctx, c, req, limit)
}

func (c *Client) GetCardComment(ctx context.Context, cardNumber int, commentID string) (*Comment, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s", c.AccountBaseURL, cardNumber, commentID)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get card comment request: %w", err)
	}

	var response Comment
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PostCardComment(ctx context.Context, cardNumber int, body string) (*Comment, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments", c.AccountBaseURL, cardNumber)

	payload := map[string]map[string]string{
		"comment": {"body": body},
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create post card comment request: %w", err)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}

	// API returns 201 Created with Location header, not a response body
	// Return a minimal Comment object with empty fields
	// The comment ID would be in the Location header but we'll use it to fetch full details
	return &Comment{}, nil
}

func (c *Client) PutCardComment(ctx context.Context, cardNumber int, commentID string, body string) (*Comment, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s", c.AccountBaseURL, cardNumber, commentID)

	payload := map[string]map[string]string{
		"comment": {"body": body},
	}

	req, err := c.newRequest(ctx, http.MethodPut, endpointURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create put card comment request: %w", err)
	}

	var response Comment
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteCardComment(ctx context.Context, cardNumber int, commentID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s", c.AccountBaseURL, cardNumber, commentID)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete card comment request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}
