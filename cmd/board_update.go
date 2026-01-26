package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var boardUpdateCmd = &cobra.Command{
	Use:   "update <board_id>",
	Short: "Update a board",
	Long:  `Update board settings such as name, access permissions, and auto-postpone period`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateBoard(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateBoard(cmd *cobra.Command, boardID string) error {
	// Check that at least one flag was explicitly set
	if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("all-access") &&
		!cmd.Flags().Changed("auto-postpone-period") && !cmd.Flags().Changed("description") {
		return fmt.Errorf("at least one flag must be provided (--name, --all-access, --auto-postpone-period, or --description)")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	// Build payload only with flags that were explicitly set
	payload := api.UpdateBoardPayload{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		payload.Name = name
	}
	if cmd.Flags().Changed("all-access") {
		allAccess, _ := cmd.Flags().GetBool("all-access")
		payload.AllAccess = &allAccess
	}
	if cmd.Flags().Changed("auto-postpone-period") {
		autoPostponePeriod, _ := cmd.Flags().GetInt("auto-postpone-period")
		payload.AutoPostponePeriod = &autoPostponePeriod
	}
	if cmd.Flags().Changed("description") {
		publicDescription, _ := cmd.Flags().GetString("description")
		payload.PublicDescription = publicDescription
	}

	err := a.Client.PutBoard(context.Background(), boardID, payload)
	if err != nil {
		return fmt.Errorf("updating board: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' updated successfully\n", boardID)
	return nil
}

func init() {
	boardUpdateCmd.Flags().StringP("name", "n", "", "Board name")
	boardUpdateCmd.Flags().Bool("all-access", false, "Allow all access to the board")
	boardUpdateCmd.Flags().Int("auto-postpone-period", 0, "Auto postpone period in days")
	boardUpdateCmd.Flags().String("description", "", "Public description of the board")

	boardCmd.AddCommand(boardUpdateCmd)
}
