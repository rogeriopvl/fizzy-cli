package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardListWithPaginationMultiplePages(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if r.URL.Path != "/cards" {
			t.Errorf("expected /cards, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
				{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Card{
				{ID: "card-3", Number: 3, Title: "Card 3", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{})
	if err != nil {
		t.Fatalf("GetCards failed: %v", err)
	}

	if len(cards) != 3 {
		t.Errorf("expected 3 cards across pages, got %d", len(cards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests for pagination, got %d", requestCount)
	}
	if cards[0].ID != "card-1" || cards[2].ID != "card-3" {
		t.Errorf("unexpected card order")
	}
}

func TestCardListWithLimitFlag(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
				{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
				{ID: "card-3", Number: 3, Title: "Card 3", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Card{
				{ID: "card-4", Number: 4, Title: "Card 4", Status: "published"},
				{ID: "card-5", Number: 5, Title: "Card 5", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{Limit: 4})
	if err != nil {
		t.Fatalf("GetCards with limit failed: %v", err)
	}

	if len(cards) != 4 {
		t.Errorf("expected 4 cards with limit=4, got %d", len(cards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests (limit not reached on page 1), got %d", requestCount)
	}
}

func TestCardListWithLimitExactBoundary(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
				{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{Limit: 2})
	if err != nil {
		t.Fatalf("GetCards with exact boundary limit failed: %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("expected 2 cards with limit=2, got %d", len(cards))
	}
	if requestCount != 1 {
		t.Errorf("expected 1 request (limit reached on page 1), got %d", requestCount)
	}
}

func TestBoardListWithPaginationMultiplePages(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if r.URL.Path != "/boards" {
			t.Errorf("expected /boards, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/boards?page=2>; rel="next"`)
			response := []api.Board{
				{ID: "board-1", Name: "Board 1", AllAccess: true},
				{ID: "board-2", Name: "Board 2", AllAccess: true},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Board{
				{ID: "board-3", Name: "Board 3", AllAccess: true},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	boards, err := client.GetBoards(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetBoards failed: %v", err)
	}

	if len(boards) != 3 {
		t.Errorf("expected 3 boards across pages, got %d", len(boards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests for pagination, got %d", requestCount)
	}
}

func TestBoardListWithLimit(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/boards?page=2>; rel="next"`)
			response := []api.Board{
				{ID: "board-1", Name: "Board 1", AllAccess: true},
				{ID: "board-2", Name: "Board 2", AllAccess: true},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Board{
				{ID: "board-3", Name: "Board 3", AllAccess: true},
				{ID: "board-4", Name: "Board 4", AllAccess: true},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	boards, err := client.GetBoards(context.Background(), &api.ListOptions{Limit: 3})
	if err != nil {
		t.Fatalf("GetBoards with limit failed: %v", err)
	}

	if len(boards) != 3 {
		t.Errorf("expected 3 boards with limit=3, got %d", len(boards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests (limit not reached on page 1), got %d", requestCount)
	}
}

func TestCardListWithPaginationAndFilters(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			boardIDs := r.URL.Query()["board_ids[]"]
			if len(boardIDs) == 0 || boardIDs[0] != "board-123" {
				t.Errorf("expected board_ids[]=board-123 in query, got %v", boardIDs)
			}

			tagIDs := r.URL.Query()["tag_ids[]"]
			if len(tagIDs) == 0 || tagIDs[0] != "tag-456" {
				t.Errorf("expected tag_ids[]=tag-456 in query, got %v", tagIDs)
			}
		}

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?board_ids%5B%5D=board-123&tag_ids%5B%5D=tag-456&page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Card{
				{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	filters := &api.CardFilters{
		BoardIDs: []string{"board-123"},
		TagIDs:   []string{"tag-456"},
	}
	cards, err := client.GetCards(context.Background(), filters)
	if err != nil {
		t.Fatalf("GetCards with pagination and filters failed: %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("expected 2 cards with filters and pagination, got %d", len(cards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests, got %d", requestCount)
	}
}

func TestPaginationStopsWhenNoNextLink(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		response := []api.Card{
			{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
			{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{})
	if err != nil {
		t.Fatalf("GetCards failed: %v", err)
	}

	if len(cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(cards))
	}
	if requestCount != 1 {
		t.Errorf("expected 1 request (no next link), got %d", requestCount)
	}
}

func TestPaginationWithEmptyPages(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			response := []api.Card{}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{})
	if err != nil {
		t.Fatalf("GetCards with empty page failed: %v", err)
	}

	if len(cards) != 1 {
		t.Errorf("expected 1 card (empty page should be handled), got %d", len(cards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests (should fetch empty page), got %d", requestCount)
	}
}

func TestPaginationErrorOnSecondPage(t *testing.T) {
	requestCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		w.Header().Set("Content-Type", "application/json")

		if requestCount == 1 {
			w.Header().Set("Link", `<`+serverURL+`/cards?page=2>; rel="next"`)
			response := []api.Card{
				{ID: "card-1", Number: 1, Title: "Card 1", Status: "published"},
				{ID: "card-2", Number: 2, Title: "Card 2", Status: "published"},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
	}))
	defer server.Close()
	serverURL = server.URL

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	cards, err := client.GetCards(context.Background(), &api.CardFilters{})

	if err == nil {
		t.Fatal("expected error on second page, got nil")
	}
	if cards != nil {
		t.Errorf("expected nil cards on error, got %d cards", len(cards))
	}
	if requestCount != 2 {
		t.Errorf("expected 2 requests before error, got %d", requestCount)
	}
}
