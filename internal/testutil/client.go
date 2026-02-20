// Package testutil
package testutil

import (
	"net/http"
	"time"

	fizzy "github.com/rogeriopvl/fizzy-go"
)

func NewTestClient(baseURL, accountSlug, boardID, accessToken string) *fizzy.Client {
	// fizzy-go requires accountSlug, use a dummy if empty
	if accountSlug == "" {
		accountSlug = "/test-account"
	}
	// fizzy-go requires accessToken, use a dummy if empty
	if accessToken == "" {
		accessToken = "test-token"
	}

	opts := []fizzy.ClientOption{
		fizzy.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
		fizzy.WithBaseURL(baseURL),
	}

	if boardID != "" {
		opts = append(opts, fizzy.WithBoard(boardID))
	}

	client, _ := fizzy.NewClient(accountSlug, accessToken, opts...)
	return client
}
