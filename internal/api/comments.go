package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetCardComments(ctx context.Context, cardNumber int) ([]Comment, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get card comments request: %w", err)
	}

	var response []Comment
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

	var response Comment
	_, err = c.decodeResponse(req, &response, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	return &response, nil
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
