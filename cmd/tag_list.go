package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	Long:  `Retrieve and display all tags in the account`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListTags(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListTags(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	tags, err := a.Client.GetTags(context.Background())
	if err != nil {
		return fmt.Errorf("fetching tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found")
		return nil
	}

	for _, tag := range tags {
		fmt.Printf("%s\n", tag.Title)
	}

	return nil
}

func init() {
	tagCmd.AddCommand(tagListCmd)
}
