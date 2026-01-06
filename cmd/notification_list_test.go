package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestNotificationListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/notifications" {
			t.Errorf("expected /notifications, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []api.Notification{
			{
				ID:        "notif-123",
				Read:      false,
				ReadAt:    "",
				CreatedAt: "2025-01-01T00:00:00Z",
				Title:     "Plain text mentions",
				Body:      "Assigned to self",
				Creator: api.User{
					ID:        "user-123",
					Name:      "David Heinemeier Hansson",
					Email:     "david@example.com",
					Role:      "owner",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.401Z",
					URL:       "http://fizzy.localhost:3006/897362094/users/03f5v9zjw7pz8717a4no1h8a7",
				},
				Card: api.CardReference{
					ID:     "card-123",
					Title:  "Plain text mentions",
					Status: "published",
					URL:    "http://fizzy.localhost:3006/897362094/cards/3",
				},
				URL: "http://fizzy.localhost:3006/897362094/notifications/03f5va03bpuvkcjemcxl73ho2",
			},
			{
				ID:        "notif-456",
				Read:      true,
				ReadAt:    "2025-01-02T00:00:00Z",
				CreatedAt: "2025-01-01T12:00:00Z",
				Title:     "Comment reply",
				Body:      "Someone replied to your comment",
				Creator: api.User{
					ID:        "user-456",
					Name:      "Jason Fried",
					Email:     "jason@example.com",
					Role:      "member",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.419Z",
					URL:       "http://fizzy.localhost:3006/897362094/users/03f5v9zjysoy0fqs9yg0ei3hq",
				},
				Card: api.CardReference{
					ID:     "card-456",
					Title:  "Fix bug",
					Status: "in_progress",
					URL:    "http://fizzy.localhost:3006/897362094/cards/4",
				},
				URL: "http://fizzy.localhost:3006/897362094/notifications/notif-456",
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

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListNotifications(cmd); err != nil {
		t.Fatalf("handleListNotifications failed: %v", err)
	}
}

func TestNotificationListCommandNoNotifications(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.Notification{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListNotifications(cmd); err != nil {
		t.Fatalf("handleListNotifications failed: %v", err)
	}
}

func TestNotificationListCommandAPIError(t *testing.T) {
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

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListNotifications(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching notifications: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestNotificationListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListNotifications(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestNotificationListCommandWithReadFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := []api.Notification{
			{
				ID:        "notif-123",
				Read:      false,
				CreatedAt: "2025-01-01T00:00:00Z",
				Title:     "Unread notification",
				Body:      "This should be filtered out",
				Creator: api.User{
					ID:        "user-123",
					Name:      "David Heinemeier Hansson",
					Email:     "david@example.com",
					Role:      "owner",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.401Z",
				},
				Card: api.CardReference{
					ID:     "card-123",
					Title:  "Test card",
					Status: "published",
				},
			},
			{
				ID:        "notif-456",
				Read:      true,
				ReadAt:    "2025-01-02T00:00:00Z",
				CreatedAt: "2025-01-01T12:00:00Z",
				Title:     "Read notification",
				Body:      "This should be shown",
				Creator: api.User{
					ID:        "user-456",
					Name:      "Jason Fried",
					Email:     "jason@example.com",
					Role:      "member",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.419Z",
				},
				Card: api.CardReference{
					ID:     "card-456",
					Title:  "Another card",
					Status: "in_progress",
				},
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

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.Flags().Set("read", "true")

	if err := handleListNotifications(cmd); err != nil {
		t.Fatalf("handleListNotifications failed: %v", err)
	}
}

func TestNotificationListCommandWithUnreadFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := []api.Notification{
			{
				ID:        "notif-123",
				Read:      false,
				CreatedAt: "2025-01-01T00:00:00Z",
				Title:     "Unread notification",
				Body:      "This should be shown",
				Creator: api.User{
					ID:        "user-123",
					Name:      "David Heinemeier Hansson",
					Email:     "david@example.com",
					Role:      "owner",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.401Z",
				},
				Card: api.CardReference{
					ID:     "card-123",
					Title:  "Test card",
					Status: "published",
				},
			},
			{
				ID:        "notif-456",
				Read:      true,
				ReadAt:    "2025-01-02T00:00:00Z",
				CreatedAt: "2025-01-01T12:00:00Z",
				Title:     "Read notification",
				Body:      "This should be filtered out",
				Creator: api.User{
					ID:        "user-456",
					Name:      "Jason Fried",
					Email:     "jason@example.com",
					Role:      "member",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.419Z",
				},
				Card: api.CardReference{
					ID:     "card-456",
					Title:  "Another card",
					Status: "in_progress",
				},
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

	cmd := notificationListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.Flags().Set("unread", "true")

	if err := handleListNotifications(cmd); err != nil {
		t.Fatalf("handleListNotifications failed: %v", err)
	}
}
