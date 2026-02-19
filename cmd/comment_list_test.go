package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCommentListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/comments" {
			t.Errorf("expected /test-account/cards/123/comments, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []fizzy.Comment{
			{
				ID:        "comment-123",
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-01T00:00:00Z",
				Creator:   fizzy.User{ID: "user-1", Name: "John Doe"},
			},
			{
				ID:        "comment-456",
				CreatedAt: "2025-01-02T00:00:00Z",
				UpdatedAt: "2025-01-02T00:00:00Z",
				Creator:   fizzy.User{ID: "user-2", Name: "Jane Doe"},
			},
		}
		// Set the body field manually since it's a struct
		response[0].Body.PlainText = "First comment"
		response[0].Body.HTML = "<p>First comment</p>"
		response[1].Body.PlainText = "Second comment"
		response[1].Body.HTML = "<p>Second comment</p>"
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := commentListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListComments(cmd, "123"); err != nil {
		t.Fatalf("handleListComments failed: %v", err)
	}
}

func TestCommentListCommandNoComments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Comment{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := commentListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListComments(cmd, "123"); err != nil {
		t.Fatalf("handleListComments failed: %v", err)
	}
}

func TestCommentListCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := commentListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListComments(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCommentListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := commentListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListComments(cmd, "123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCommentListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := commentListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListComments(cmd, "123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching comments: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}
