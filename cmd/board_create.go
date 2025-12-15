package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var (
	boardName               string
	boardAllAccess          bool
	boardAutoPostponePeriod int
	boardPublicDescription  string
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

	payload := api.CreateBoardPayload{
		Name:               boardName,
		AllAccess:          boardAllAccess,
		AutoPostponePeriod: boardAutoPostponePeriod,
		PublicDescription:  boardPublicDescription,
	}

	_, err := a.Client.PostBoards(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating board: %w", err)
	}

	fmt.Printf("✓ Board '%s' created successfully\n", boardName)
	return nil
}

func init() {
	boardCreateCmd.Flags().StringVarP(&boardName, "name", "n", "", "Board name (required)")
	boardCreateCmd.MarkFlagRequired("name")
	boardCreateCmd.Flags().BoolVar(&boardAllAccess, "all-access", false, "Allow all access to the board")
	boardCreateCmd.Flags().IntVar(&boardAutoPostponePeriod, "auto-postpone-period", 0, "Auto postpone period in days")
	boardCreateCmd.Flags().StringVar(&boardPublicDescription, "description", "", "Public description of the board")

	boardCmd.AddCommand(boardCreateCmd)
}
