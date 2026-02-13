package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilgehannal/cursor-config/curset/internal/collection"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local .cursor operations",
	Long:  "Commands for inspecting the local .cursor folder.",
}

var localListCmd = &cobra.Command{
	Use:   "list",
	Short: "List local .cursor contents as a collection",
	Long:  "Scans the current directory's .cursor/ folder and displays its contents in collection.json format.",
	Run: func(cmd *cobra.Command, args []string) {
		cursorDir := ".cursor"

		info, err := os.Stat(cursorDir)
		if err != nil || !info.IsDir() {
			fmt.Fprintln(os.Stderr, "Error: .cursor/ directory not found in current directory")
			os.Exit(1)
		}

		col := make(collection.Collection)

		// Read top-level directories under .cursor/ (e.g. rules, commands)
		entries, err := os.ReadDir(cursorDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading .cursor/: %v\n", err)
			os.Exit(1)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			objType := entry.Name()
			objPath := filepath.Join(cursorDir, objType)

			children, err := os.ReadDir(objPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", objPath, err)
				continue
			}

			var items []string
			for _, child := range children {
				if child.IsDir() {
					items = append(items, child.Name())
				} else {
					// Use filename without extension
					name := child.Name()
					ext := filepath.Ext(name)
					items = append(items, strings.TrimSuffix(name, ext))
				}
			}

			if len(items) > 0 {
				col[objType] = items
			}
		}

		cf := &collection.CollectionFile{
			Collections: map[string]collection.Collection{
				"local": col,
			},
		}

		jsonStr, err := cf.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(jsonStr)
	},
}

func init() {
	localCmd.AddCommand(localListCmd)
}
