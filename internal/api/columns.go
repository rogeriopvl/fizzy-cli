package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetColumns(ctx context.Context) ([]Column, error) {
	if c.BoardBaseURL == "" {
		return nil, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get columns request: %w", err)
	}

	var response []Column
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetColumn(ctx context.Context, columnID string) (*Column, error) {
	if c.BoardBaseURL == "" {
		return nil, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns/" + columnID

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response Column
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PostColumns(ctx context.Context, payload CreateColumnPayload) (bool, error) {
	if c.BoardBaseURL == "" {
		return false, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns"

	body := map[string]CreateColumnPayload{"column": payload}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create column request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusCreated)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PutColumn(ctx context.Context, columnID string, payload UpdateColumnPayload) error {
	if c.BoardBaseURL == "" {
		return fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns/" + columnID

	body := map[string]UpdateColumnPayload{"column": payload}

	req, err := c.newRequest(ctx, http.MethodPut, endpointURL, body)
	if err != nil {
		return fmt.Errorf("failed to create update column request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteColumn(ctx context.Context, columnID string) error {
	if c.BoardBaseURL == "" {
		return fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns/" + columnID

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete column request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}
