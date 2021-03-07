package cli

import (
	"os"

	"github.com/gustavooferreira/wcrawler/internal/graph"
	"github.com/spf13/cobra"
)

func newViewCmd() *cobra.Command {
	var (
		inputFilePath  string
		outputFilePath string
		noautoopen     bool
	)

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View web links relationships in the browser",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			iFile, err := os.Open(inputFilePath)
			if err != nil {
				return err
			}
			defer iFile.Close()

			oFile, err := os.Create(outputFilePath)
			if err != nil {
				return err
			}

			defer oFile.Close()

			v := graph.NewViewer(iFile, oFile)
			err = v.Run()
			if err != nil {
				return nil
			}

			if !noautoopen {
				graph.Openbrowser(outputFilePath)
			}

			return err
		},
	}

	viewCmd.Flags().StringVarP(&inputFilePath, "input", "i", "./web_graph.json", "file containing the data")
	viewCmd.Flags().StringVarP(&outputFilePath, "output", "o", "./web_graph.html", "HTML output file")
	viewCmd.Flags().BoolVarP(&noautoopen, "noautoopen", "n", false, "don't open browser automatically")

	return viewCmd
}
