package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestColumnDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/columns/col-456" {
			t.Errorf("expected /test-account/boards/board-123/columns/col-456, got %s", r.URL.Path)
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

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteColumn(cmd, "col-456"); err != nil {
		t.Fatalf("handleDeleteColumn failed: %v", err)
	}
}

func TestColumnDeleteCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Column not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteColumn(cmd, "nonexistent-col")
	if err == nil {
		t.Errorf("expected error for column not found")
	}
	if err.Error() != "deleting column: unexpected status code 404: Column not found" {
		t.Errorf("expected column not found error, got %v", err)
	}
}

func TestColumnDeleteCommandForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You don't have permission to delete this column"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error for forbidden access")
	}
	if err.Error() != "deleting column: unexpected status code 403: You don't have permission to delete this column" {
		t.Errorf("expected permission error, got %v", err)
	}
}

func TestColumnDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting column: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestColumnDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := columnDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
