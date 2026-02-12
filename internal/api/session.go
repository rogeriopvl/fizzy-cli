package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) DeleteSession(ctx context.Context) error {
	endpointURL := c.BaseURL + "/session"

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}
