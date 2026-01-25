package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetCards(ctx context.Context, filters CardFilters) ([]Card, error) {
	endpointURL := c.AccountBaseURL + "/cards"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get cards request: %w", err)
	}

	if len(filters.BoardIDs) > 0 {
		q := req.URL.Query()
		for _, boardID := range filters.BoardIDs {
			q.Add("board_ids[]", boardID)
		}
		req.URL.RawQuery = q.Encode()
	}

	var response []Card
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetCard(ctx context.Context, cardNumber int) (*Card, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get card by id request: %w", err)
	}

	var response Card
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PostCards(ctx context.Context, payload CreateCardPayload) (bool, error) {
	if c.BoardBaseURL == "" {
		return false, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/cards"

	body := map[string]CreateCardPayload{"card": payload}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create card request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusCreated)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PutCard(ctx context.Context, cardNumber int, payload UpdateCardPayload) (*Card, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d", c.AccountBaseURL, cardNumber)

	body := map[string]UpdateCardPayload{"card": payload}

	req, err := c.newRequest(ctx, http.MethodPut, endpointURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create update card request: %w", err)
	}

	var response Card
	_, err = c.decodeResponse(req, &response, http.StatusOK)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteCard(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete card request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardsClosure(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/closure", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create closure card request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardNotNow(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/not_now", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create post not now request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardTriage(ctx context.Context, cardNumber int, columnID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/triage", c.AccountBaseURL, cardNumber)

	body := map[string]any{"column_id": columnID}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create post triage request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) DeleteCardTriage(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/triage", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete triage request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardWatch(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/watch", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create post watch request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) DeleteCardWatch(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/watch", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete watch request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardGoldenness(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/goldness", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create post goldness request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) DeleteCardGoldenness(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/goldness", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete goldness request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) DeleteCardsClosure(ctx context.Context, cardNumber int) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/closure", c.AccountBaseURL, cardNumber)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete closure card request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardAssignments(ctx context.Context, cardNumber int, userID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/assignments", c.AccountBaseURL, cardNumber)

	body := map[string]string{"assignee_id": userID}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create assignment request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostCardTagging(ctx context.Context, cardNumber int, tagTitle string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/cards/%d/taggings", c.AccountBaseURL, cardNumber)

	body := map[string]string{"tag_title": tagTitle}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create tagging request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}
