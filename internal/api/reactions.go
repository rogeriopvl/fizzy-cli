package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetCommentReactions(ctx context.Context, cardNumber int, commentID string) ([]Reaction, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s/reactions", c.AccountBaseURL, cardNumber, commentID)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get comment reactions request: %w", err)
	}

	var response []Reaction
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) PostCommentReaction(ctx context.Context, cardNumber int, commentID string, content string) (*Reaction, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s/reactions", c.AccountBaseURL, cardNumber, commentID)

	payload := map[string]map[string]string{
		"reaction": {"content": content},
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create post comment reaction request: %w", err)
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
	// Return a Reaction object with the content field set for reference
	return &Reaction{Content: content}, nil
}

func (c *Client) DeleteCommentReaction(ctx context.Context, cardNumber int, commentID string, reactionID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/comments/%s/reactions/%s", c.AccountBaseURL, cardNumber, commentID, reactionID)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete comment reaction request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}
