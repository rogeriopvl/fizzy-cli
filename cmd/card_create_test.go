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
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board-123/cards" {
			t.Errorf("expected /boards/board-123/cards, got %s", r.URL.Path)
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
		var payload map[string]api.CreateCardPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		cardPayload := payload["card"]
		if cardPayload.Title != "Implement feature" {
			t.Errorf("expected title 'Implement feature', got %s", cardPayload.Title)
		}
		if cardPayload.Description != "A detailed description" {
			t.Errorf("expected description 'A detailed description', got %s", cardPayload.Description)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cardTitle = "Implement feature"
	cardDescription = "A detailed description"
	cardStatus = ""
	cardImageURL = ""
	cardTagIDs = []string{}
	cardCreatedAt = ""
	cardLastActiveAt = ""

	if err := handleCreateCard(cmd); err != nil {
		t.Fatalf("handleCreateCard failed: %v", err)
	}
}

func TestCardCreateCommandWithAllFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.CreateCardPayload
		json.Unmarshal(body, &payload)

		cardPayload := payload["card"]
		if cardPayload.Title != "Fix bug" {
			t.Errorf("expected title 'Fix bug', got %s", cardPayload.Title)
		}
		if cardPayload.Status != "in_progress" {
			t.Errorf("expected status 'in_progress', got %s", cardPayload.Status)
		}
		if cardPayload.ImageURL != "https://example.com/image.jpg" {
			t.Errorf("expected image URL 'https://example.com/image.jpg', got %s", cardPayload.ImageURL)
		}
		if len(cardPayload.TagIDS) != 2 || cardPayload.TagIDS[0] != "tag-1" || cardPayload.TagIDS[1] != "tag-2" {
			t.Errorf("expected tag IDs [tag-1, tag-2], got %v", cardPayload.TagIDS)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cardTitle = "Fix bug"
	cardDescription = ""
	cardStatus = "in_progress"
	cardImageURL = "https://example.com/image.jpg"
	cardTagIDs = []string{"tag-1", "tag-2"}
	cardCreatedAt = ""
	cardLastActiveAt = ""

	if err := handleCreateCard(cmd); err != nil {
		t.Fatalf("handleCreateCard failed: %v", err)
	}
}

func TestCardCreateCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: ""},
	}

	cmd := cardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cardTitle = "Test"
	cardDescription = ""
	cardStatus = ""
	cardImageURL = ""
	cardTagIDs = []string{}
	cardCreatedAt = ""
	cardLastActiveAt = ""

	err := handleCreateCard(cmd)
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "no board selected" {
		t.Errorf("expected 'no board selected' error, got %v", err)
	}
}

func TestCardCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cardTitle = "Test"
	cardDescription = ""
	cardStatus = ""
	cardImageURL = ""
	cardTagIDs = []string{}
	cardCreatedAt = ""
	cardLastActiveAt = ""

	err := handleCreateCard(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardCreateCommandAPIError(t *testing.T) {
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

	cmd := cardCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cardTitle = "Test"
	cardDescription = ""
	cardStatus = ""
	cardImageURL = ""
	cardTagIDs = []string{}
	cardCreatedAt = ""
	cardLastActiveAt = ""

	err := handleCreateCard(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating card: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}
