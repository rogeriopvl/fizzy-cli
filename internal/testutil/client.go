// Package testutil
package testutil

import (
	"net/http"
	"time"

	"github.com/rogeriopvl/fizzy/internal/api"
)

func NewTestClient(baseURL, accountSlug, boardID, accessToken string) *api.Client {
	accountBaseURL := baseURL + accountSlug

	var boardURL string
	if boardID != "" {
		boardURL = accountBaseURL + "/boards" + "/" + boardID
	}

	return &api.Client{
		BaseURL:        baseURL,
		AccountBaseURL: accountBaseURL,
		BoardBaseURL:   boardURL,
		AccessToken:    accessToken,
		HTTPClient:     &http.Client{Timeout: 30 * time.Second},
	}
}
