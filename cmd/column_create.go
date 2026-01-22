package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/colors"
	"github.com/spf13/cobra"
)

func buildColorAliases() map[string]api.Color {
	aliases := make(map[string]api.Color)
	for _, colorDef := range colors.All {
		aliases[strings.ToLower(colorDef.Name)] = api.Color(colorDef.CSSValue)
	}
	return aliases
}

func getAvailableColors() string {
	var names []string
	for _, colorDef := range colors.All {
		names = append(names, strings.ToLower(colorDef.Name))
	}
	return strings.Join(names, ", ")
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

	// Read flag values directly from command
	name, _ := cmd.Flags().GetString("name")
	colorStr, _ := cmd.Flags().GetString("color")

	payload := api.CreateColumnPayload{
		Name: name,
	}

	if colorStr != "" {
		colorAliases := buildColorAliases()
		color, ok := colorAliases[colorStr]
		if !ok {
			return fmt.Errorf("invalid color '%s'. Available colors: %s", colorStr, getAvailableColors())
		}
		payload.Color = &color
	}

	_, err := a.Client.PostColumns(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating column: %w", err)
	}

	fmt.Printf("✓ Column '%s' created successfully\n", name)
	return nil
}

func init() {
	columnCreateCmd.Flags().StringP("name", "n", "", "Column name (required)")
	columnCreateCmd.MarkFlagRequired("name")
	columnCreateCmd.Flags().String("color", "", fmt.Sprintf("Column color (optional). Available: %s", getAvailableColors()))

	columnCmd.AddCommand(columnCreateCmd)
}
