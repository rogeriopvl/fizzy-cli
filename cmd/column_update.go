package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var columnUpdateCmd = &cobra.Command{
	Use:   "update <column_id>",
	Short: "Update a column",
	Long:  `Update column settings such as name and color`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateColumn(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateColumn(cmd *cobra.Command, columnID string) error {
	// Check that at least one flag was explicitly set
	if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("color") {
		return fmt.Errorf("at least one flag must be provided (--name or --color)")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	// Build payload only with flags that were explicitly set
	payload := fizzy.UpdateColumnPayload{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		payload.Name = name
	}

	if cmd.Flags().Changed("color") {
		colorStr, _ := cmd.Flags().GetString("color")
		colorAliases := buildColorAliases()
		color, ok := colorAliases[colorStr]
		if !ok {
			return fmt.Errorf("invalid color '%s'. Available colors: %s", colorStr, getAvailableColors())
		}
		payload.Color = &color
	}

	err := a.Client.UpdateColumn(context.Background(), columnID, payload)
	if err != nil {
		return fmt.Errorf("updating column: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Column '%s' updated successfully\n", columnID)
	return nil
}

func init() {
	columnUpdateCmd.Flags().StringP("name", "n", "", "Column name")
	columnUpdateCmd.Flags().String("color", "", fmt.Sprintf("Column color (optional). Available: %s", getAvailableColors()))

	columnCmd.AddCommand(columnUpdateCmd)
}
