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
	"github.com/spf13/cobra"
)

func TestActivityListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/activities" {
			t.Errorf("expected /test-account/activities, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Activity{
			{ID: "evt-1", Action: "card_closed", CreatedAt: "2026-03-25T15:11:04Z", Description: "X closed Y"},
			{ID: "evt-2", Action: "comment_created", CreatedAt: "2026-03-25T14:17:22Z", Description: "X commented on Y"},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := activityListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListActivities(cmd); err != nil {
		t.Fatalf("handleListActivities failed: %v", err)
	}
}

func TestActivityListCommandFiltersAndLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		creators := q["creator_ids[]"]
		if len(creators) != 1 || creators[0] != "user-1" {
			t.Errorf("expected creator_ids[]=user-1, got %v", creators)
		}
		boards := q["board_ids[]"]
		if len(boards) != 2 || boards[0] != "board-A" || boards[1] != "board-B" {
			t.Errorf("expected board_ids[]=board-A,board-B, got %v", boards)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Activity{
			{ID: "evt-1", Action: "card_closed", CreatedAt: "2026-03-25T15:11:04Z", Description: "X"},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := activityListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--creator", "user-1", "--board", "board-A,board-B", "--limit", "5"})

	if err := handleListActivities(cmd); err != nil {
		t.Fatalf("handleListActivities failed: %v", err)
	}
}

func TestActivityListCommandEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Activity{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{}
	cmd.Flags().StringSlice("creator", nil, "")
	cmd.Flags().StringSlice("board", nil, "")
	cmd.Flags().Int("limit", 0, "")
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListActivities(cmd); err != nil {
		t.Fatalf("handleListActivities failed: %v", err)
	}
}

func TestActivityListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := activityListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListActivities(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching activities: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestActivityListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := activityListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListActivities(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
