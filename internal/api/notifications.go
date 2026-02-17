package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) GetNotifications(ctx context.Context, opts *ListOptions) ([]Notification, error) {
	endpointURL := c.AccountBaseURL + "/notifications"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get notifications request: %w", err)
	}

	limit := 0
	if opts != nil {
		limit = opts.Limit
	}

	return fetchAllPages[Notification](ctx, c, req, limit)
}

func (c *Client) GetNotification(ctx context.Context, notificationID string) (*Notification, error) {
	endpointURL := fmt.Sprintf("%s/notifications/%s", c.AccountBaseURL, notificationID)

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get notification request: %w", err)
	}

	var response Notification
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) PostNotificationReading(ctx context.Context, notificationID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/notifications/%s/reading", c.AccountBaseURL, notificationID)

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create mark notification as read request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) DeleteNotificationReading(ctx context.Context, notificationID string) (bool, error) {
	endpointURL := fmt.Sprintf("%s/notifications/%s/reading", c.AccountBaseURL, notificationID)

	req, err := c.newRequest(ctx, http.MethodDelete, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete notification request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) PostBulkNotificationsReading(ctx context.Context) (bool, error) {
	endpointURL := c.AccountBaseURL + "/notifications/bulk_reading"

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create bulk notifications reading request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusNoContent)
	if err != nil {
		return false, err
	}

	return true, nil
}
