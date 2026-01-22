package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestColumnCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board-123/columns" {
			t.Errorf("expected /boards/board-123/columns, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.CreateColumnPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		columnPayload := payload["column"]
		if columnPayload.Name != "Todo" {
			t.Errorf("expected name 'Todo', got %s", columnPayload.Name)
		}
		if columnPayload.Color != nil {
			t.Errorf("expected no color, got %v", columnPayload.Color)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", "/columns/col-123")
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Todo"})

	if err := handleCreateColumn(cmd); err != nil {
		t.Fatalf("handleCreateColumn failed: %v", err)
	}
}

func TestColumnCreateCommandWithColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.CreateColumnPayload
		json.Unmarshal(body, &payload)

		columnPayload := payload["column"]
		if columnPayload.Name != "In Progress" {
			t.Errorf("expected name 'In Progress', got %s", columnPayload.Name)
		}
		if columnPayload.Color == nil {
			t.Error("expected color to be set")
		}
		if *columnPayload.Color != "var(--color-card-4)" {
			t.Errorf("expected color 'var(--color-card-4)', got %s", *columnPayload.Color)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "In Progress", "--color", "lime"})

	if err := handleCreateColumn(cmd); err != nil {
		t.Fatalf("handleCreateColumn failed: %v", err)
	}
}

func TestColumnCreateCommandInvalidColor(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Todo", "--color", "invalid"})

	err := handleCreateColumn(cmd)
	if err == nil {
		t.Errorf("expected error for invalid color")
	}
	if err.Error() != "invalid color 'invalid'. Available colors: blue, gray, tan, yellow, lime, aqua, violet, purple, pink" {
		t.Errorf("expected invalid color error, got %v", err)
	}
}

func TestColumnCreateCommandNoBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Todo", "--color", ""})

	err := handleCreateColumn(cmd)
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "creating column: please select a board first with 'fizzy use --board <board_name>'" {
		t.Errorf("expected 'board not selected' error, got %v", err)
	}
}

func TestColumnCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Todo", "--color", ""})

	err := handleCreateColumn(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestColumnCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Todo", "--color", ""})

	err := handleCreateColumn(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating column: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}
