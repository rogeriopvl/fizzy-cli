package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var boardCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new board",
	Long:  `Create a new board in Fizzy`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateBoard(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateBoard(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	// Read flag values directly from command
	name, _ := cmd.Flags().GetString("name")
	allAccess, _ := cmd.Flags().GetBool("all-access")
	autoPostponePeriod, _ := cmd.Flags().GetInt("auto-postpone-period")
	publicDescription, _ := cmd.Flags().GetString("description")

	payload := api.CreateBoardPayload{
		Name:               name,
		AllAccess:          allAccess,
		AutoPostponePeriod: autoPostponePeriod,
		PublicDescription:  publicDescription,
	}

	_, err := a.Client.PostBoards(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating board: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Board '%s' created successfully\n", name)
	return nil
}

func init() {
	boardCreateCmd.Flags().StringP("name", "n", "", "Board name (required)")
	boardCreateCmd.MarkFlagRequired("name")
	boardCreateCmd.Flags().Bool("all-access", false, "Allow all access to the board")
	boardCreateCmd.Flags().Int("auto-postpone-period", 0, "Auto postpone period in days")
	boardCreateCmd.Flags().String("description", "", "Public description of the board")

	boardCmd.AddCommand(boardCreateCmd)
}
