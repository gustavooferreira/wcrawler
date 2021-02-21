package cli

import (
	"net/http"
	"time"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/spf13/cobra"
)

func newExploreCmd() *cobra.Command {
	var (
		file    string
		stats   bool
		workers uint
		timeout uint
		depth   uint
		client  *http.Client
	)

	exploreCmd := &cobra.Command{
		Use:   "explore URL",
		Short: "Explore the web by following links up to a pre-determined depth",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]

			// Validate url is actually a url
			// fmt.Printf("file: %+v - workers: %+v - timeout: %+v - stats: %+v - url: %+v\n", file, workers, timeout, stats, url)

			client = &http.Client{
				Timeout: time.Second * time.Duration(timeout),
			}

			c := core.NewCrawler(client, file, stats, workers, depth)
			c.Run(url)
		},
	}

	exploreCmd.Flags().StringVarP(&file, "file", "f", "./web_graph.json", "file to save results")
	exploreCmd.Flags().BoolVarP(&stats, "stats", "s", true, "show live stats")
	exploreCmd.Flags().UintVarP(&workers, "workers", "w", 10, "number of workers making concurrent requests")
	exploreCmd.Flags().UintVarP(&timeout, "timeout", "t", 10, "HTTP requests timeout")
	exploreCmd.Flags().UintVarP(&depth, "depth", "d", 10, "depth of recursion")

	return exploreCmd
}
