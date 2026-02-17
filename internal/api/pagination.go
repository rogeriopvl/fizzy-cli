package api

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

func parseNextLink(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	// Match pattern: <URL>; rel="next"
	re := regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)
	matches := re.FindStringSubmatch(linkHeader)
	if len(matches) >= 2 {
		return matches[1]
	}

	// Also try without quotes: <URL>; rel=next
	re = regexp.MustCompile(`<([^>]+)>;\s*rel=next`)
	matches = re.FindStringSubmatch(linkHeader)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

type ListOptions struct {
	Limit int
}

// fetchAllPages handles pagination by iterating through all pages and returning all items.
// Since Go doesn't support generic methods on concrete types, callers must unmarshal into
// the concrete type before calling this function. This is handled by decodeResponse.
func fetchAllPages[T any](ctx context.Context, client *Client, req *http.Request, limit int) ([]T, error) {
	var allItems []T

	for {
		var pageItems []T
		resp, err := client.decodeResponse(req, &pageItems)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, pageItems...)

		// Check if we've reached the limit
		if limit > 0 && len(allItems) >= limit {
			allItems = allItems[:limit]
			break
		}

		// Check if there are more pages
		if resp.NextURL == "" {
			break
		}

		// Create request for next page
		req, err = client.newRequest(ctx, http.MethodGet, resp.NextURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create next page request: %w", err)
		}
	}

	return allItems, nil
}
