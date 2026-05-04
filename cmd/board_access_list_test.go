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

func TestBoardAccessListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/accesses" {
			t.Errorf("expected /test-account/boards/board-123/accesses, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.BoardAccesses{
			BoardID:   "board-123",
			AllAccess: true,
			Users: []fizzy.BoardAccess{
				{User: fizzy.User{ID: "u-1", Name: "Alice"}, HasAccess: true, Involvement: "creator"},
				{User: fizzy.User{ID: "u-2", Name: "Bob"}, HasAccess: false},
			},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardAccessListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListBoardAccesses(cmd, "board-123"); err != nil {
		t.Fatalf("handleListBoardAccesses failed: %v", err)
	}
}

func TestBoardAccessListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardAccessListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListBoardAccesses(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching board accesses: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardAccessListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardAccessListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListBoardAccesses(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
