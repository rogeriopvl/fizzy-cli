package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetTags(ctx context.Context) ([]Tag, error) {
	endpointURL := c.AccountBaseURL + "/tags"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get tags request: %w", err)
	}

	var response []Tag
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
