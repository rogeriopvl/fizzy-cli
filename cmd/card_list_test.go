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
	"github.com/spf13/cobra"
)

func newCardListCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().StringSliceP("tag", "t", []string{}, "Filter by tag ID (can be used multiple times)")
	cmd.Flags().StringSliceP("assignee", "a", []string{}, "Filter by assignee user ID (can be used multiple times)")
	cmd.Flags().StringSlice("creator", []string{}, "Filter by creator user ID (can be used multiple times)")
	cmd.Flags().StringSlice("closer", []string{}, "Filter by closer user ID (can be used multiple times)")
	cmd.Flags().StringSlice("card", []string{}, "Filter to specific card ID (can be used multiple times)")
	cmd.Flags().String("indexed-by", "", "Filter by status: all, closed, not_now, stalled, postponing_soon, golden")
	cmd.Flags().String("sorted-by", "", "Sort order: latest, newest, oldest")
	cmd.Flags().BoolP("unassigned", "u", false, "Show only unassigned cards")
	cmd.Flags().String("created-in", "", "Filter by creation date")
	cmd.Flags().String("closed-in", "", "Filter by closure date")
	cmd.Flags().StringSliceP("search", "s", []string{}, "Search terms (can be used multiple times)")
	return cmd
}

func TestCardListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards" {
			t.Errorf("expected /test-account/cards, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		boardIDs := r.URL.Query()["board_ids[]"]
		if len(boardIDs) == 0 || boardIDs[0] != "board-123" {
			t.Errorf("expected board_ids[]=board-123 in query, got %v", boardIDs)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []fizzy.Card{
			{
				ID:        "card-123",
				Number:    1,
				Title:     "Implement feature",
				Status:    "in_progress",
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			{
				ID:        "card-456",
				Number:    2,
				Title:     "Fix bug",
				Status:    "todo",
				CreatedAt: "2025-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandNoCards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching cards: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardListCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: ""},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "no board selected" {
		t.Errorf("expected 'no board selected' error, got %v", err)
	}
}

func TestCardListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardListCommandWithTagFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tagIDs := r.URL.Query()["tag_ids[]"]
		if len(tagIDs) != 2 || tagIDs[0] != "tag-123" || tagIDs[1] != "tag-456" {
			t.Errorf("expected tag_ids[]=tag-123&tag_ids[]=tag-456, got %v", tagIDs)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--tag", "tag-123", "--tag", "tag-456"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithAssigneeFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assigneeIDs := r.URL.Query()["assignee_ids[]"]
		if len(assigneeIDs) != 1 || assigneeIDs[0] != "user-123" {
			t.Errorf("expected assignee_ids[]=user-123, got %v", assigneeIDs)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--assignee", "user-123"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithIndexedByFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		indexedBy := r.URL.Query().Get("indexed_by")
		if indexedBy != "closed" {
			t.Errorf("expected indexed_by=closed, got %s", indexedBy)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--indexed-by", "closed"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithSortedByFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortedBy := r.URL.Query().Get("sorted_by")
		if sortedBy != "newest" {
			t.Errorf("expected sorted_by=newest, got %s", sortedBy)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--sorted-by", "newest"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithUnassignedFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assignmentStatus := r.URL.Query().Get("assignment_status")
		if assignmentStatus != "unassigned" {
			t.Errorf("expected assignment_status=unassigned, got %s", assignmentStatus)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--unassigned"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithSearchFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		terms := r.URL.Query()["terms[]"]
		if len(terms) != 2 || terms[0] != "bug" || terms[1] != "critical" {
			t.Errorf("expected terms[]=bug&terms[]=critical, got %v", terms)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--search", "bug", "--search", "critical"})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandWithMultipleFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		boardIDs := r.URL.Query()["board_ids[]"]
		if len(boardIDs) == 0 || boardIDs[0] != "board-123" {
			t.Errorf("expected board_ids[]=board-123, got %v", boardIDs)
		}

		tagIDs := r.URL.Query()["tag_ids[]"]
		if len(tagIDs) != 1 || tagIDs[0] != "tag-123" {
			t.Errorf("expected tag_ids[]=tag-123, got %v", tagIDs)
		}

		assigneeIDs := r.URL.Query()["assignee_ids[]"]
		if len(assigneeIDs) != 1 || assigneeIDs[0] != "user-456" {
			t.Errorf("expected assignee_ids[]=user-456, got %v", assigneeIDs)
		}

		sortedBy := r.URL.Query().Get("sorted_by")
		if sortedBy != "latest" {
			t.Errorf("expected sorted_by=latest, got %s", sortedBy)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := newCardListCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--tag", "tag-123",
		"--assignee", "user-456",
		"--sorted-by", "latest",
	})

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}
