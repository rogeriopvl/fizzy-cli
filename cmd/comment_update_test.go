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
)

func TestCommentUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/comments/comment-456" {
			t.Errorf("expected /test-account/cards/123/comments/comment-456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		commentPayload := payload["comment"]
		if commentPayload["body"] != "Updated comment text" {
			t.Errorf("expected body 'Updated comment text', got %s", commentPayload["body"])
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Comment{
			ID:        "comment-456",
			CreatedAt: "2025-01-01T00:00:00Z",
			UpdatedAt: "2025-01-02T00:00:00Z",
			Creator:   fizzy.User{ID: "user-1", Name: "John Doe"},
		}
		response.Body.PlainText = "Updated comment text"
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := commentUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--body", "Updated comment text",
	})

	if err := handleUpdateComment(cmd, "123", "comment-456"); err != nil {
		t.Fatalf("handleUpdateComment failed: %v", err)
	}
}

func TestCommentUpdateCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := commentUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--body", "Updated text",
	})

	err := handleUpdateComment(cmd, "not-a-number", "comment-456")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCommentUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := commentUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--body", "Updated text",
	})

	err := handleUpdateComment(cmd, "123", "comment-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCommentUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Comment not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := commentUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--body", "Updated text",
	})

	err := handleUpdateComment(cmd, "123", "comment-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating comment: unexpected status code 404: Comment not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
