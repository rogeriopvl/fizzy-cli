// Package api
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rogeriopvl/fizzy/internal/colors"
)

type Client struct {
	BaseURL        string
	AccountBaseURL string
	BoardBaseURL   string
	AccessToken    string
	HTTPClient     *http.Client
}

func NewClient(accountSlug string, boardID string) (*Client, error) {
	baseURL := "https://app.fizzy.do"
	accountBaseURL := baseURL + accountSlug

	var boardBaseURL string
	if boardID != "" {
		boardBaseURL = accountBaseURL + "/boards" + "/" + boardID
	}

	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return nil, fmt.Errorf("FIZZY_ACCESS_TOKEN environment variable is not set")
	}

	return &Client{
		BaseURL:        baseURL,
		AccountBaseURL: accountBaseURL,
		BoardBaseURL:   boardBaseURL,
		AccessToken:    token,
		HTTPClient:     &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// newRequest makes an HTTP request with the required headers setup
func (c *Client) newRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// decodeResponse executes a request and decodes the JSON response into v
// If expectedStatus is 0, it defaults to http.StatusOK
// If v is nil, the response body is not decoded
func (c *Client) decodeResponse(req *http.Request, v any, expectedStatus ...int) (int, error) {
	expectedCode := http.StatusOK
	if len(expectedStatus) > 0 {
		expectedCode = expectedStatus[0]
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != expectedCode {
		body, _ := io.ReadAll(res.Body)
		return 0, fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(res.Body).Decode(v); err != nil {
			return 0, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return res.StatusCode, nil
}

func (c *Client) GetBoards(ctx context.Context) ([]Board, error) {
	endpointURL := c.AccountBaseURL + "/boards"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response []Board
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) PostBoards(ctx context.Context, payload CreateBoardPayload) (bool, error) {
	endpointURL := c.AccountBaseURL + "/boards"

	body := map[string]CreateBoardPayload{"board": payload}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create board request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusCreated)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) GetColumns(ctx context.Context) ([]Column, error) {
	if c.BoardBaseURL == "" {
		return nil, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns"

	req, err := c.newRequest(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get columns request: %w", err)
	}

	var response []Column
	_, err = c.decodeResponse(req, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) PostColumns(ctx context.Context, payload CreateColumnPayload) (bool, error) {
	if c.BoardBaseURL == "" {
		return false, fmt.Errorf("please select a board first with 'fizzy use --board <board_name>'")
	}

	endpointURL := c.BoardBaseURL + "/columns"

	body := map[string]CreateColumnPayload{"column": payload}

	req, err := c.newRequest(ctx, http.MethodPost, endpointURL, body)
	if err != nil {
		return false, fmt.Errorf("failed to create column request: %w", err)
	}

	_, err = c.decodeResponse(req, nil, http.StatusCreated)
	if err != nil {
		return false, err
	}

	return true, nil
}

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

func (c *Client) GetCard(ctx context.Context, cardNumber string) (*Card, error) {
	endpointURL := fmt.Sprintf("%s/cards/%s", c.AccountBaseURL, cardNumber)

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

type Board struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AllAccess bool   `json:"all_access"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
	Creator   User   `json:"creator"`
}

type CreateBoardPayload struct {
	Name               string `json:"name"`
	AllAccess          bool   `json:"all_access"`
	AutoPostponePeriod int    `json:"auto_postpone_period"`
	PublicDescription  string `json:"public_description"`
}

type Column struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Color     ColorObject `json:"color"`
	CreatedAt string      `json:"created_at"`
}

type ColorObject struct {
	Name  string `json:"name"`
	Value Color  `json:"value"`
}

type CreateColumnPayload struct {
	Name  string `json:"name"`
	Color *Color `json:"color,omitempty"`
}

type Card struct {
	ID              string   `json:"id"`
	Number          int      `json:"number"`
	Title           string   `json:"title"`
	Status          string   `json:"status"`
	Description     string   `json:"description"`
	DescriptionHTML string   `json:"description_html"`
	ImageURL        string   `json:"image_url"`
	Tags            []string `json:"tags"`
	Golden          bool     `json:"golden"`
	LastActiveAt    string   `json:"last_active_at"`
	CreatedAt       string   `json:"created_at"`
	URL             string   `json:"url"`
	Board           Board    `json:"board"`
	Creator         User     `json:"creator"`
	CommentsURL     string   `json:"comments_url"`
}

type CardFilters struct {
	BoardIDs         []string
	TagIDs           []string
	AssigneeIDs      []string
	CreatorIDs       []string
	CloserIDs        []string
	CardIDs          []string
	IndexedBy        string
	SortedBy         string
	AssignmentStatus string
	CreationStatus   string
	ClosureStatus    string
	Terms            []string
}

type CreateCardPayload struct {
	Title        string   `json:"title"`
	Description  string   `json:"description,omitempty"`
	Status       string   `json:"status,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
	TagIDS       []string `json:"tag_ids,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	LastActiveAt string   `json:"last_active_at,omitempty"`
}

type GetMyIdentityResponse struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      User   `json:"user"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email_address"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
}

type Color string

// Color constants using centralized definitions
var (
	Blue   Color = Color(colors.Blue.CSSValue)
	Gray   Color = Color(colors.Gray.CSSValue)
	Tan    Color = Color(colors.Tan.CSSValue)
	Yellow Color = Color(colors.Yellow.CSSValue)
	Lime   Color = Color(colors.Lime.CSSValue)
	Aqua   Color = Color(colors.Aqua.CSSValue)
	Violet Color = Color(colors.Violet.CSSValue)
	Purple Color = Color(colors.Purple.CSSValue)
	Pink   Color = Color(colors.Pink.CSSValue)
)
