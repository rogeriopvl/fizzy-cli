package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestColumnCardsCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/columns/col-1/cards" {
			t.Errorf("expected /test-account/boards/board-123/columns/col-1/cards, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{
			{ID: "c-1", Number: 1, Title: "First"},
			{ID: "c-2", Number: 2, Title: "Second"},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCardsCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListColumnCards(cmd, "col-1"); err != nil {
		t.Fatalf("handleListColumnCards failed: %v", err)
	}
}

func TestColumnCardsCommandEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCardsCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListColumnCards(cmd, "col-1"); err != nil {
		t.Fatalf("handleListColumnCards failed: %v", err)
	}
}

func TestColumnCardsCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCardsCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumnCards(cmd, "col-1")
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "fetching column cards: no board selected: use SetBoard or WithBoard when creating the client" {
		t.Errorf("expected board error, got %v", err)
	}
}

func TestColumnCardsCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCardsCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumnCards(cmd, "col-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching column cards: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestColumnCardsCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := columnCardsCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumnCards(cmd, "col-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
