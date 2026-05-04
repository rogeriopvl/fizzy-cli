package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var boardEntropyCmd = &cobra.Command{
	Use:   "entropy <board_id>",
	Short: "Update a board's auto-postpone period",
	Long:  `Update the auto-postpone period (in days) for a specific board. Requires board admin permission.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleBoardEntropy(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleBoardEntropy(cmd *cobra.Command, boardID string) error {
	if !cmd.Flags().Changed("auto-postpone-days") {
		return fmt.Errorf("--auto-postpone-days is required")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	days, _ := cmd.Flags().GetInt("auto-postpone-days")
	payload := fizzy.EntropyPayload{AutoPostponePeriodInDays: days}

	board, err := a.Client.UpdateBoardEntropy(context.Background(), boardID, payload)
	if err != nil {
		return fmt.Errorf("updating board entropy: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' auto-postpone period set to %d days\n", board.Name, board.AutoPostponePeriodInDays)
	return nil
}

func init() {
	boardEntropyCmd.Flags().Int("auto-postpone-days", 0, "Auto-postpone period in days (required)")
	boardCmd.AddCommand(boardEntropyCmd)
}
