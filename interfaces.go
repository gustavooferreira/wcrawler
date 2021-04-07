package wcrawler

import (
	"time"
)

// Connector describes the connector interface
type Connector interface {
	GetLinks(rawURL string) (statusCode int, links []URLEntity, latency time.Duration, err error)
}

type StatsManager interface {
	// Functional options pattern
	UpdateStats(updates ...func(StatsManager))
	SetAppState(state AppState)

	RunOutputFlusher()
}
