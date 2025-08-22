package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arbor",
	Short: "Arbor - A Git repository analysis tool",
	Long: `Arbor is a CLI tool for analyzing Git repositories.
It provides various commands to analyze code metrics, commit history, and generate visualizations.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands here
	rootCmd.AddCommand(locCmd)
}
