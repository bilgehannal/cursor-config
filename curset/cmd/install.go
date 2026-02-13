package cmd

import (
	"fmt"
	"os"

	"github.com/bilgehannal/cursor-config/curset/internal/collection"
	"github.com/bilgehannal/cursor-config/curset/internal/github"
	"github.com/bilgehannal/cursor-config/curset/internal/installer"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [collection-name]",
	Short: "Install a collection",
	Long:  "Installs a named collection into the current directory's .cursor/ folder.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
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

		col, ok := cf.Collections[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: collection '%s' not found\n", name)
			fmt.Fprintln(os.Stderr, "\nAvailable collections:")
			for _, n := range cf.SortedNames() {
				fmt.Fprintf(os.Stderr, "  - %s\n", n)
			}
			os.Exit(1)
		}

		inst, err := installer.NewInstaller(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := inst.Install(col, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
