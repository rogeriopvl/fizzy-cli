package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestTagListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tags" {
			t.Errorf("expected /tags, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []api.Tag{
			{
				ID:        "tag-123",
				Title:     "bug",
				CreatedAt: "2025-01-01T00:00:00Z",
				URL:       "http://fizzy.localhost:3006/897362094/cards?tag_ids[]=tag-123",
			},
			{
				ID:        "tag-456",
				Title:     "feature",
				CreatedAt: "2025-01-02T00:00:00Z",
				URL:       "http://fizzy.localhost:3006/897362094/cards?tag_ids[]=tag-456",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := tagListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListTags(cmd); err != nil {
		t.Fatalf("handleListTags failed: %v", err)
	}
}

func TestTagListCommandNoTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := tagListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListTags(cmd); err != nil {
		t.Fatalf("handleListTags with no tags failed: %v", err)
	}
}

func TestTagListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := tagListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListTags(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}

func TestTagListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := tagListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListTags(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
