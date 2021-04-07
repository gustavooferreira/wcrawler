package wcrawler

import (
	"time"
)

// Connector describes the connector interface.
type Connector interface {
	GetLinks(rawURL string) (statusCode int, links []URLEntity, latency time.Duration, err error)
}

// StatsManager represents a tracker of statistics related to the crawler.
// This interface is unfortunately quite big as it needs to support several
// operations on the statistics it keeps track of.
type StatsManager interface {
	SetAppState(state AppState)
	SetLinksInQueue(value int)
	IncDecLinksInQueue(value int)
	SetLinksCount(value int)
	IncDecLinksCount(value int)
	SetErrorsCount(value int)
	IncDecErrorsCount(value int)
	SetWorkersRunning(value int)
	IncDecWorkersRunning(value int)
	SetTotalRequestsCount(value int)
	IncDecTotalRequestsCount(value int)
	SetDepth(value int)
	IncDecDepth(value int)
	AddLatencySample(value time.Duration)
	RunOutputFlusher()
}
