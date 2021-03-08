package wcrawler

import (
	"time"

	"github.com/gustavooferreira/wcrawler/internal/stats"
)

// Connector describes the connector interface
type Connector interface {
	GetLinks(rawURL string) (statusCode int, links []URLEntity, latency time.Duration, err error)
}

type StatsManager interface {
	UpdateStats(updates ...func(*stats.StatsCLIOutput))
	RunOutputFlusher()
}
