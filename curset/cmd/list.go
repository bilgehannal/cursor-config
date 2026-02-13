package cmd

import (
	"fmt"
	"os"

	"github.com/bilgehannal/cursor-config/curset/internal/collection"
	"github.com/bilgehannal/cursor-config/curset/internal/github"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available collections",
	Long:  "Fetches the collection.json from the remote repository and displays all available collections.",
	Run: func(cmd *cobra.Command, args []string) {
		client := github.NewClient()

		data, err := client.FetchCollectionJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		cf, err := collection.Parse(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		cf.PrintTables()
	},
}
