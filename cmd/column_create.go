package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var (
	columnName  string
	columnColor string
)

var colorAliases = map[string]api.Color{
	"blue":   api.Blue,
	"gray":   api.Gray,
	"tan":    api.Tan,
	"yellow": api.Yellow,
	"lime":   api.Lime,
	"aqua":   api.Aqua,
	"violet": api.Violet,
	"purple": api.Purple,
	"pink":   api.Pink,
}

var columnCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new column",
	Long:  `Create a new column in the selected board. Color is optional.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateColumn(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateColumn(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	payload := api.CreateColumnPayload{
		Name: columnName,
	}

	if columnColor != "" {
		color, ok := colorAliases[columnColor]
		if !ok {
			return fmt.Errorf("invalid color '%s'. Available colors: blue, gray, tan, yellow, lime, aqua, violet, purple, pink", columnColor)
		}
		payload.Color = &color
	}

	_, err := a.Client.PostColumns(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating column: %w", err)
	}

	fmt.Printf("✓ Column '%s' created successfully\n", columnName)
	return nil
}

func init() {
	columnCreateCmd.Flags().StringVarP(&columnName, "name", "n", "", "Column name (required)")
	columnCreateCmd.MarkFlagRequired("name")
	columnCreateCmd.Flags().StringVar(&columnColor, "color", "", "Column color (optional). Available: blue, gray, tan, yellow, lime, aqua, violet, purple, pink")

	columnCmd.AddCommand(columnCreateCmd)
}
