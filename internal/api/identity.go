package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetMyIdentity(ctx context.Context) (*GetMyIdentityResponse, error) {
	endpointURL := c.BaseURL + "/my/identity"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response GetMyIdentityResponse
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
