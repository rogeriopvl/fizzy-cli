package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetCardStep(ctx context.Context, cardNumber int, stepID string) (*Step, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/steps/%s", c.AccountBaseURL, cardNumber, stepID)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get card step request: %w", err)
	}

	var response Step
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PostCardStep(ctx context.Context, cardNumber int, content string, completed bool) (*Step, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/steps", c.AccountBaseURL, cardNumber)

	payload := map[string]map[string]any{
		"step": {
			"content":   content,
			"completed": completed,
		},
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create post card step request: %w", err)
	}

	var response Step
	_, err = c.decodeResponse(req, &response, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PutCardStep(ctx context.Context, cardNumber int, stepID string, content *string, completed *bool) (*Step, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/steps/%s", c.AccountBaseURL, cardNumber, stepID)

	stepPayload := make(map[string]any)
	if content != nil {
		stepPayload["content"] = *content
	}
	if completed != nil {
		stepPayload["completed"] = *completed
	}

	payload := map[string]map[string]any{
		"step": stepPayload,
	}

	req, err := c.newRequest(ctx, http.MethodPut, endpointURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create put card step request: %w", err)
	}

	var response Step
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteCardStep(ctx context.Context, cardNumber int, stepID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/steps/%s", c.AccountBaseURL, cardNumber, stepID)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete card step request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}
