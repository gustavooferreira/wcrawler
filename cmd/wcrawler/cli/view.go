package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newViewCmd() *cobra.Command {
	var (
		file string
	)

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View web links relationships in the browser",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("file: %+v\n", file)

			// Save the html file in the same folder where the json file is located
		},
	}

	viewCmd.Flags().StringVarP(&file, "file", "f", "./web_graph.json", "file containing the data")

	return viewCmd
}
