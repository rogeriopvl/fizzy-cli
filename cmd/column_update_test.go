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

func newTestUpdateColumnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "update <column_id>",
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().StringP("name", "n", "", "Column name")
	cmd.Flags().String("color", "", "Column color")
	return cmd
}

func TestColumnUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/columns/col-456" {
			t.Errorf("expected /test-account/boards/board-123/columns/col-456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateColumnPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		columnPayload := payload["column"]
		if columnPayload.Name != "Updated Column" {
			t.Errorf("expected name 'Updated Column', got %s", columnPayload.Name)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Column"})

	if err := handleUpdateColumn(cmd, "col-456"); err != nil {
		t.Fatalf("handleUpdateColumn failed: %v", err)
	}
}

func TestColumnUpdateCommandWithColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateColumnPayload
		json.Unmarshal(body, &payload)

		columnPayload := payload["column"]
		if columnPayload.Name != "Progress" {
			t.Errorf("expected name 'Progress', got %s", columnPayload.Name)
		}
		if columnPayload.Color == nil {
			t.Error("expected Color to be set")
		} else if *columnPayload.Color != fizzy.ColorLime {
			t.Errorf("expected color Lime, got %s", *columnPayload.Color)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Progress", "--color", "lime"})

	if err := handleUpdateColumn(cmd, "col-456"); err != nil {
		t.Fatalf("handleUpdateColumn failed: %v", err)
	}
}

func TestColumnUpdateCommandInvalidColor(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--color", "invalid-color"})

	err := handleUpdateColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error for invalid color")
	}
	errMsg := err.Error()
	if errMsg != "invalid color 'invalid-color'. Available colors: blue, gray, tan, yellow, lime, aqua, violet, purple, pink" {
		t.Errorf("expected invalid color error, got %v", err)
	}
}

func TestColumnUpdateCommandNoFlags(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "at least one flag must be provided (--name or --color)" {
		t.Errorf("expected flag requirement error, got %v", err)
	}
}

func TestColumnUpdateCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Column not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Column"})

	err := handleUpdateColumn(cmd, "nonexistent-col")
	if err == nil {
		t.Errorf("expected error for column not found")
	}
	if err.Error() != "updating column: unexpected status code 404: Column not found" {
		t.Errorf("expected column not found error, got %v", err)
	}
}

func TestColumnUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Column"})

	err := handleUpdateColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating column: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestColumnUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := newTestUpdateColumnCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Column"})

	err := handleUpdateColumn(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
