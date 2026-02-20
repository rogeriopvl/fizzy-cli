package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestBoardCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123" {
			t.Errorf("expected /test-account/boards/board-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Board{
			ID:        "board-123",
			Name:      "Test Board",
			AllAccess: true,
			CreatedAt: "2025-12-05T19:36:35.534Z",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{
			SelectedBoard: "board-123",
		},
	}

	cmd := boardCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowBoard(cmd); err != nil {
		t.Fatalf("handleShowBoard failed: %v", err)
	}
}

func TestBoardCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{
			SelectedBoard: "",
		},
	}

	cmd := boardCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowBoard(cmd)
	if err == nil {
		t.Errorf("expected error when no board selected")
	}
	if err.Error() != "no board selected" {
		t.Errorf("expected 'no board selected' error, got %v", err)
	}
}

func TestBoardCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{
			SelectedBoard: "board-123",
		},
	}

	cmd := boardCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowBoard(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}
