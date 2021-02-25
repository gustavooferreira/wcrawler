package cli

import (
	"net/http"
	"time"

	"github.com/gustavooferreira/wcrawler/pkg/core"
	"github.com/spf13/cobra"
)

func newExploreCmd() *cobra.Command {
	var (
		file            string
		stats           bool
		workers         uint
		timeout         uint
		depth           uint
		stayinsubdomain bool
		client          *http.Client
	)

	exploreCmd := &cobra.Command{
		Use:   "explore URL",
		Short: "Explore the web by following links up to a pre-determined depth",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			client = &http.Client{
				Timeout: time.Second * time.Duration(timeout),
			}

			connector := core.NewWebClient(client)
			c, err := core.NewCrawler(connector, url, file, stats, stayinsubdomain, workers, depth)
			if err != nil {
				return err
			}
			c.Run()
			return nil
		},
	}

	exploreCmd.Flags().StringVarP(&file, "output", "o", "./web_graph.json", "file to save results")
	exploreCmd.Flags().BoolVarP(&stats, "stats", "s", true, "show live stats")
	exploreCmd.Flags().BoolVarP(&stayinsubdomain, "stayinsubdomain", "z", false, "follow links only in the same subdomain")
	exploreCmd.Flags().UintVarP(&workers, "workers", "w", 10, "number of workers making concurrent requests")
	exploreCmd.Flags().UintVarP(&timeout, "timeout", "t", 10, "HTTP requests timeout")
	exploreCmd.Flags().UintVarP(&depth, "depth", "d", 10, "depth of recursion")

	return exploreCmd
}
