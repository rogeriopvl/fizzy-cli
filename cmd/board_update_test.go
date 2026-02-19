package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
	"github.com/spf13/cobra"
)

func TestBoardUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123" {
			t.Errorf("expected /test-account/boards/board-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateBoardPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		boardPayload := payload["board"]
		if boardPayload.Name != "Updated Board" {
			t.Errorf("expected name 'Updated Board', got %s", boardPayload.Name)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Board"})

	if err := handleUpdateBoard(cmd, "board-123"); err != nil {
		t.Fatalf("handleUpdateBoard failed: %v", err)
	}
}

func TestBoardUpdateCommandWithAllFlags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateBoardPayload
		json.Unmarshal(body, &payload)

		boardPayload := payload["board"]
		if boardPayload.Name != "Updated Board" {
			t.Errorf("expected name 'Updated Board', got %s", boardPayload.Name)
		}
		if boardPayload.AllAccess == nil || !*boardPayload.AllAccess {
			t.Error("expected AllAccess to be pointer to true")
		}
		if boardPayload.AutoPostponePeriod == nil || *boardPayload.AutoPostponePeriod != 14 {
			t.Errorf("expected AutoPostponePeriod pointer to 14, got %v", boardPayload.AutoPostponePeriod)
		}
		if boardPayload.PublicDescription != "Updated description" {
			t.Errorf("expected description 'Updated description', got %s", boardPayload.PublicDescription)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--name", "Updated Board",
		"--all-access",
		"--auto-postpone-period", "14",
		"--description", "Updated description",
	})

	if err := handleUpdateBoard(cmd, "board-123"); err != nil {
		t.Fatalf("handleUpdateBoard failed: %v", err)
	}
}

func TestBoardUpdateCommandZeroValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateBoardPayload
		json.Unmarshal(body, &payload)

		boardPayload := payload["board"]
		// Verify that false value is sent (not omitted)
		if boardPayload.AllAccess == nil {
			t.Error("expected AllAccess to be set (pointer), but got nil")
		}
		if *boardPayload.AllAccess != false {
			t.Error("expected AllAccess to be false")
		}
		// Verify that 0 value is sent (not omitted)
		if boardPayload.AutoPostponePeriod == nil {
			t.Error("expected AutoPostponePeriod to be set (pointer), but got nil")
		}
		if *boardPayload.AutoPostponePeriod != 0 {
			t.Errorf("expected AutoPostponePeriod to be 0, got %d", *boardPayload.AutoPostponePeriod)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--all-access=false",
		"--auto-postpone-period", "0",
	})

	if err := handleUpdateBoard(cmd, "board-123"); err != nil {
		t.Fatalf("handleUpdateBoard failed: %v", err)
	}
}

func TestBoardUpdateCommandNoFlags(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	// Create a fresh command to avoid flag pollution from other tests
	cmd := &cobra.Command{
		Use:  "update <board_id>",
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String("name", "", "Board name")
	cmd.Flags().Bool("all-access", false, "Allow all access to the board")
	cmd.Flags().Int("auto-postpone-period", 0, "Auto postpone period in days")
	cmd.Flags().String("description", "", "Public description of the board")

	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "at least one flag must be provided (--name, --all-access, --auto-postpone-period, or --description)" {
		t.Errorf("expected flag requirement error, got %v", err)
	}
}

func TestBoardUpdateCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Board not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Board"})

	err := handleUpdateBoard(cmd, "nonexistent-board")
	if err == nil {
		t.Errorf("expected error for board not found")
	}
	if err.Error() != "updating board: unexpected status code 404: Board not found" {
		t.Errorf("expected board not found error, got %v", err)
	}
}

func TestBoardUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Board"})

	err := handleUpdateBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating board: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Board"})

	err := handleUpdateBoard(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
