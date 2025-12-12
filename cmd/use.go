package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the active board or account",
	Long:  `Set the active board or account to use for subsequent commands`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUse(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUse(cmd *cobra.Command) error {
	board, _ := cmd.Flags().GetString("board")
	account, _ := cmd.Flags().GetString("account")

	if board == "" && account == "" {
		return fmt.Errorf("must specify either --board or --account")
	}

	if board != "" && account != "" {
		return fmt.Errorf("cannot specify both --board and --account")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if board != "" {
		a := app.FromContext(cmd.Context())
		if a == nil || a.Client == nil {
			return fmt.Errorf("API client not available")
		}

		boards, err := a.Client.GetBoards(context.Background())
		if err != nil {
			return fmt.Errorf("fetching boards: %w", err)
		}

		var boardID string
		for _, b := range boards {
			if b.Name == board {
				boardID = b.ID
				break
			}
		}

		if boardID == "" {
			return fmt.Errorf("board '%s' not found", board)
		}

		cfg.SelectedBoard = boardID
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("Selected board: %s\n", board)
	}

	if account != "" {
		cfg.SelectedAccount = account
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("Selected account: %s\n", account)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().String("board", "", "Board name to use")
	useCmd.Flags().String("account", "", "Account slug to use")
}
