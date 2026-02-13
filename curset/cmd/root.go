package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var version = "dev"
var gitignoreFlag bool

var rootCmd = &cobra.Command{
	Use:   "curset",
	Short: "Cursor config collection manager",
	Long:  "curset is a CLI tool for managing .cursor folder configurations from curated collections.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if gitignoreFlag {
			if err := addCursorToGitignore(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to update .gitignore: %v\n", err)
			}
		}
	},
}

// Execute runs the root command.
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&gitignoreFlag, "gitignore", "g", false, "Add .cursor/ to .gitignore in the current directory")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(localCmd)
}

// addCursorToGitignore adds ".cursor/" to the current directory's .gitignore file.
// Creates the file if it doesn't exist. Skips if ".cursor/" is already present.
func addCursorToGitignore() error {
	const gitignorePath = ".gitignore"
	const cursorEntry = ".cursor/"

	// Check if .gitignore exists and already contains .cursor/
	if data, err := os.ReadFile(gitignorePath); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == cursorEntry {
				fmt.Println(".gitignore: .cursor/ already present")
				return nil
			}
		}
	}

	// Append .cursor/ to .gitignore (create if not exists)
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if file ends with a newline, add one if needed
	info, err := f.Stat()
	if err != nil {
		return err
	}

	prefix := ""
	if info.Size() > 0 {
		// Read last byte to check for trailing newline
		content, err := os.ReadFile(gitignorePath)
		if err != nil {
			return err
		}
		if len(content) > 0 && content[len(content)-1] != '\n' {
			prefix = "\n"
		}
	}

	if _, err := fmt.Fprintf(f, "%s%s\n", prefix, cursorEntry); err != nil {
		return err
	}

	fmt.Println(".gitignore: added .cursor/")
	return nil
}
