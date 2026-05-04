package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var boardAccessListCmd = &cobra.Command{
	Use:   "list <board_id>",
	Short: "List user access for a board",
	Long:  `Retrieve and display the access list for a board`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListBoardAccesses(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListBoardAccesses(cmd *cobra.Command, boardID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	opts := &fizzy.ListOptions{}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		opts.Limit = limit
	}

	accesses, err := a.Client.GetBoardAccesses(context.Background(), boardID, opts)
	if err != nil {
		return fmt.Errorf("fetching board accesses: %w", err)
	}

	return ui.DisplayBoardAccesses(cmd.OutOrStdout(), accesses)
}

func init() {
	boardAccessListCmd.Flags().IntP("limit", "l", 0, "Maximum number of users to return (0 = no limit)")
	boardAccessCmd.AddCommand(boardAccessListCmd)
}
