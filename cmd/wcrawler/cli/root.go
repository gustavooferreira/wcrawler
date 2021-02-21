package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "wcrawler",
		Short:        "A cli tool for crawling the web",
		SilenceUsage: true,
		// SilenceErrors: true,
	}

	// Init sub commands
	exploreCmd := newExploreCmd()
	viewCmd := newViewCmd()

	rootCmd.AddCommand(exploreCmd, viewCmd)
	return rootCmd
}
