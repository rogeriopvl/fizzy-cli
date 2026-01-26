package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestBoardDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board-123" {
			t.Errorf("expected /boards/board-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteBoard(cmd, "board-123"); err != nil {
		t.Fatalf("handleDeleteBoard failed: %v", err)
	}
}

func TestBoardDeleteCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Board not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteBoard(cmd, "nonexistent-board")
	if err == nil {
		t.Errorf("expected error for board not found")
	}
	if err.Error() != "deleting board: unexpected status code 404: Board not found" {
		t.Errorf("expected board not found error, got %v", err)
	}
}

func TestBoardDeleteCommandForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You don't have permission to delete this board"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for forbidden access")
	}
	if err.Error() != "deleting board: unexpected status code 403: You don't have permission to delete this board" {
		t.Errorf("expected permission error, got %v", err)
	}
}

func TestBoardDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting board: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
