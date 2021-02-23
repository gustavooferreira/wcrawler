package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newViewCmd() *cobra.Command {
	var (
		input  string
		output string
	)

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View web links relationships in the browser",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("file: %+v | %+v\n", input, output)

			// Save the html file in the same folder where the json file is located
		},
	}

	viewCmd.Flags().StringVarP(&input, "input", "i", "./web_graph.json", "file containing the data")
	viewCmd.Flags().StringVarP(&output, "output", "o", "./web_graph.html", "HTML output")

	return viewCmd
}
