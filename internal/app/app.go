// Package app
package app

import (
	"context"
	"fmt"
	"os"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/rogeriopvl/fizzy-cli/internal/config"
)

type App struct {
	Client *api.Client
	Config *config.Config
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return &App{Config: cfg}, nil // No token set, app will handle gracefully
	}

	client, err := api.NewClient(cfg.SelectedAccount, cfg.SelectedBoard)
	if err != nil {
		return nil, fmt.Errorf("creating API client: %w", err)
	}

	return &App{Client: client, Config: cfg}, nil
}

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const appContextKey contextKey = "app"

func FromContext(ctx context.Context) *App {
	app, _ := ctx.Value(appContextKey).(*App)
	return app
}

func (a *App) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appContextKey, a)
}
