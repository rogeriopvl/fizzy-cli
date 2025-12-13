// Package testutil
package testutil

import "github.com/rogeriopvl/fizzy-cli/internal/api"

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
	}
}
