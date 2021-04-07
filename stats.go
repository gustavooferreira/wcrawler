package wcrawler

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gustavooferreira/wcrawler/internal/ring"
)

// TODO: we can do better here.
// I'm acquiring a lock every single time, and sometimes these stats's functions
// are called together, wasting time acquiring and releasing a lock for each call.
// I had originally implemented this using the "functional options pattern" and it
// worked great. But it's not easy to abstract it away with an interface.

// StatsCLIOutWriter keeps track of stats and writes to a writer up to date stats.
type StatsCLIOutWriter struct {

	// keep a reference to where to print stats
	writer io.Writer

	// Read only vars
	// --------------
	// enables or disables the presentation of errors in the output
	showErrorsFlag bool
	// total
	totalWorkersCount int
	// max depth level provided by user
	maxDepthLevel int

	// mu protects access to the fields below
	mu sync.Mutex

	// Crawler state
	state AppState

	// This is the total number of links still to be checked
	// This number will keep increasing as new links are found.
	linksInQueue int
	linksCount   int
	errorCounts  int
	// number of workers running currently
	workersRunning int
	// number of HTTP requests made
	totalRequestsCount int
	// current level of depth
	depth int

	// latency
	lMin float64
	lMax float64
	// Use these two numbers to compute average
	lAvgSum   float64
	lAvgCount float64

	// Requests per second
	rps int

	// List of errors that happen during crawling
	errorsList ring.Buffer
}

// NewStatsCLIOutWriter returns a new StatsCLIOutWriter.
func NewStatsCLIOutWriter(writer io.Writer, showErrors bool, totalWorkersCount int, depth int) *StatsCLIOutWriter {
	sm := StatsCLIOutWriter{
		writer:            writer,
		showErrorsFlag:    showErrors,
		state:             AppState_IDLE,
		totalWorkersCount: totalWorkersCount,
		maxDepthLevel:     depth,
		errorsList:        ring.New(10)}
	return &sm
}

func (sm *StatsCLIOutWriter) SetAppState(state AppState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state = state
}

func (sm *StatsCLIOutWriter) SetLinksInQueue(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.linksInQueue = value
}

func (sm *StatsCLIOutWriter) IncDecLinksInQueue(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.linksInQueue += value
}

func (sm *StatsCLIOutWriter) SetLinksCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.linksInQueue = value
}

func (sm *StatsCLIOutWriter) IncDecLinksCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.linksCount += value
}

func (sm *StatsCLIOutWriter) SetErrorsCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.errorCounts = value
}

func (sm *StatsCLIOutWriter) IncDecErrorsCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.errorCounts += value
}

func (sm *StatsCLIOutWriter) SetWorkersRunning(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.workersRunning = value
}

func (sm *StatsCLIOutWriter) IncDecWorkersRunning(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.workersRunning += value
}

func (sm *StatsCLIOutWriter) SetTotalRequestsCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.totalRequestsCount = value
}

func (sm *StatsCLIOutWriter) IncDecTotalRequestsCount(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.totalRequestsCount += value
}

func (sm *StatsCLIOutWriter) SetDepth(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.depth = value
}

func (sm *StatsCLIOutWriter) IncDecDepth(value int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.depth += value
}

func (sm *StatsCLIOutWriter) AddLatencySample(value time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	valueF := value.Seconds()
	if sm.lAvgCount == 0 {
		sm.lMax = valueF
		sm.lMin = valueF
		sm.lAvgCount = 1
		sm.lAvgSum = valueF
	} else {
		if valueF > sm.lMax {
			sm.lMax = valueF
		}

		if valueF < sm.lMin {
			sm.lMin = valueF
		}

		sm.lAvgCount++
		sm.lAvgSum += valueF
	}
}

func (sm *StatsCLIOutWriter) AddErrorEntry(value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.errorsList.Add(value)
}

// This functions writes the updated stats to an io.Writer
// Run this in a goroutine
func (sm *StatsCLIOutWriter) RunOutputFlusher() {

	fmtStr := "Crawler State: %11s\n" +
		"Links in Cache: %10d     Depth:       (%3d/%3d)\n" +
		"Links in Queue: %10d     Workers:   (%4d/%4d)\n" +
		"Total Req Count: %9d     Errors: %5d (%5.2f%%)\n" +
		"Requests/s: %14d\n" +
		"Latency --------------- (in  seconds) ---------------\n" +
		"Min: %6.3f    -     Avg: %6.3f     -    Max: %6.3f\n"

	errorsStr := "Last 10 errors max ----------------------------------\n"

	// If zero samples, don't display latency

	// Setup counter (up to 5 times) since we wait 200 milliseconds,
	// when we get to 5 it means it's time to look at the totalrequest count
	// and subtract from the previous value

	previousTotalRequestCount := 0
	cyclesCount := 0
	var rps int = 0

	for {
		sm.mu.Lock()

		errorsPerc := 0.0
		if sm.totalRequestsCount != 0 {
			errorsPerc = 100 * float64(sm.errorCounts) / float64(sm.totalRequestsCount)
		}

		// requests per second
		// This is a rough estimation
		if cyclesCount == 5 {
			cyclesCount = 0
			rps = sm.totalRequestsCount - previousTotalRequestCount
			previousTotalRequestCount = sm.totalRequestsCount
		}

		// truncate error messages to 50 chars
		// last 10 errors:
		// - error 1
		// - error 2
		// - error etc

		// if len(sm.errors) != 0
		// append string at the end with the errors formatted
		// Only show last 10 errors
		// Errors (last 10):
		// - one error per line
		// - two errors, etc

		var statsBuf strings.Builder

		fmt.Fprintf(&statsBuf, fmtStr, sm.state, sm.linksCount, sm.depth, sm.maxDepthLevel,
			sm.linksInQueue, sm.workersRunning, sm.totalWorkersCount, sm.totalRequestsCount,
			sm.errorCounts, errorsPerc, rps, sm.lMin, sm.lAvgSum/sm.lAvgCount, sm.lMax)

		if sm.showErrorsFlag {
			if sm.errorsList.Len() != 0 {
				fmt.Fprintf(&statsBuf, errorsStr)

				for _, item := range sm.errorsList.ReadAll() {
					fmt.Fprintf(&statsBuf, "- %s\n", item)
				}
			}
		}

		fmt.Fprint(sm.writer, statsBuf.String())
		sm.mu.Unlock()

		// Only stop when AppState == Finished!
		if sm.state == AppState_Finished {
			break
		}

		time.Sleep(time.Millisecond * 200)
		cyclesCount++
	}
}
