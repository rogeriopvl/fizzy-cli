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

func TestBoardCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards" {
			t.Errorf("expected /boards, got %s", r.URL.Path)
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
		var payload map[string]api.CreateBoardPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		boardPayload := payload["board"]
		if boardPayload.Name != "Test Board" {
			t.Errorf("expected name 'Test Board', got %s", boardPayload.Name)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", "/boards/board-123")
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Test Board"})

	if err := handleCreateBoard(cmd); err != nil {
		t.Fatalf("handleCreateBoard failed: %v", err)
	}
}

func TestBoardCreateCommandWithAllFlags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.CreateBoardPayload
		json.Unmarshal(body, &payload)

		boardPayload := payload["board"]
		if !boardPayload.AllAccess {
			t.Error("expected AllAccess to be true")
		}
		if boardPayload.AutoPostponePeriod != 7 {
			t.Errorf("expected AutoPostponePeriod 7, got %d", boardPayload.AutoPostponePeriod)
		}
		if boardPayload.PublicDescription != "Team project" {
			t.Errorf("expected description 'Team project', got %s", boardPayload.PublicDescription)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--name", "Test Board",
		"--all-access",
		"--auto-postpone-period", "7",
		"--description", "Team project",
	})

	if err := handleCreateBoard(cmd); err != nil {
		t.Fatalf("handleCreateBoard failed: %v", err)
	}
}

func TestBoardCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Test Board"})

	err := handleCreateBoard(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating board: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Test Board"})

	err := handleCreateBoard(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
