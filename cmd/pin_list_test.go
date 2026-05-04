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

func TestPinListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/my/pins" {
			t.Errorf("expected /test-account/my/pins, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{
			{ID: "card-1", Number: 1, Title: "First"},
			{ID: "card-2", Number: 2, Title: "Second"},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := pinListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListPins(cmd); err != nil {
		t.Fatalf("handleListPins failed: %v", err)
	}
}

func TestPinListCommandNoPins(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := pinListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListPins(cmd); err != nil {
		t.Fatalf("handleListPins failed: %v", err)
	}
}

func TestPinListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := pinListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListPins(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching pinned cards: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestPinListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := pinListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListPins(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
