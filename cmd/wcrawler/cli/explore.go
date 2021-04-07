package cli

import (
	"net/http"
	"os"
	"time"

	"github.com/gustavooferreira/wcrawler"
	"github.com/spf13/cobra"
)

func newExploreCmd() *cobra.Command {
	var (
		filePath        string
		nostats         bool
		workers         uint
		timeout         uint
		retry           uint
		depth           uint
		stayinsubdomain bool
		treemode        bool
		client          *http.Client
	)

	exploreCmd := &cobra.Command{
		Use:   "explore URL",
		Short: "Explore the web by following links up to a pre-determined depth",
		Long: "Explore the web by following links up to a pre-determined depth.\n" +
			"A depth of zero means no limit.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			client = &http.Client{
				Timeout: time.Second * time.Duration(timeout),
			}

			f, err := os.Create(filePath)
			if err != nil {
				return err
			}

			defer f.Close()

			connector := wcrawler.NewWebClient(client)
			c, err := wcrawler.NewCrawler(connector, url, int(retry), f, !nostats, stayinsubdomain, treemode, int(workers), int(depth))
			if err != nil {
				return err
			}
			c.Run()
			return nil
		},
	}

	exploreCmd.Flags().StringVarP(&filePath, "output", "o", "./web_graph.json", "file to save results")
	exploreCmd.Flags().BoolVarP(&nostats, "nostats", "s", false, "don't show live stats")
	exploreCmd.Flags().UintVarP(&workers, "workers", "w", 100, "number of workers making concurrent requests")
	exploreCmd.Flags().UintVarP(&timeout, "timeout", "t", 10, "HTTP requests timeout in seconds")
	exploreCmd.Flags().UintVarP(&retry, "retry", "r", 2, "retry requests when they timeout")
	exploreCmd.Flags().UintVarP(&depth, "depth", "d", 5, "depth of recursion")
	exploreCmd.Flags().BoolVarP(&stayinsubdomain, "stayinsubdomain", "z", false, "follow links only in the same subdomain")
	exploreCmd.Flags().BoolVarP(&treemode, "treemode", "m", false, "doesn't add links which would point back to known nodes")

	return exploreCmd
}
