package wcrawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type StatsCLIOutput struct {

	// keep a reference to were to print stats
	writer *uilive.Writer

	// Read only vars
	// --------------
	totalWorkersCount int
	// depth level provided by user
	maxDepthLevel int

	// mu protects access to the fields below
	mu sync.Mutex

	// Crawler State
	State AppState

	// This is the total number of links still to be checked
	// This number will keep increasing as new links are found.
	linksInQueue       int
	linksCount         int
	errorCounts        int
	workersRunning     int
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
	errors [10]error
}

func NewStatsManager(totalWorkersCount int, depth int) *StatsCLIOutput {
	sm := StatsCLIOutput{State: AppState_IDLE, totalWorkersCount: totalWorkersCount, maxDepthLevel: depth}
	sm.writer = uilive.New()

	return &sm
}

// UpdateStats updates the stats.
func (sm *StatsCLIOutput) UpdateStats(updates ...func(*StatsCLIOutput)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, update := range updates {
		update(sm)
	}
}

func SetAppState(value AppState) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.State = value
	}
}

func SetLinksInQueue(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.linksInQueue = value
	}
}

func IncDecLinksInQueue(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.linksInQueue += value
	}
}

func SetLinksCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.linksCount = value
	}
}

func IncDecLinksCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.linksCount += value
	}
}

func SetErrorsCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.errorCounts = value
	}
}

func IncDecErrorsCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.errorCounts += value
	}
}

func SetWorkersRunning(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.workersRunning = value
	}
}

func IncDecWorkersRunning(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.workersRunning += value
	}
}

func SetTotalRequestsCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.totalRequestsCount = value
	}
}

func IncDecTotalRequestsCount(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.totalRequestsCount += value
	}
}

func SetDepth(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.depth = value
	}
}

func IncDecDepth(value int) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.depth += value
	}
}

func AddLatencySample(value time.Duration) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
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
}

// This functions writes the updated stats to an io.Writer
// Run this in a goroutine
func (sm *StatsCLIOutput) RunOutputFlusher() {
	sm.writer.Start()

	fmtStr := "Crawler State: %11s\n" +
		"Links in Cache: %10d     Depth:       (%3d/%3d)\n" +
		"Links in Queue: %10d     Workers:   (%4d/%4d)\n" +
		"Total Req Count: %9d     Errors: %5d (%5.2f%%)\n" +
		"Requests/s: %14d\n" +
		"Latency --------------- (in  seconds) ---------------\n" +
		"Min: %6.3f    -     Avg: %6.3f     -    Max: %6.3f\n"
		// "Last 10 errors --------------------------------------\n" +
		// "%s"

		// If zero samples, don't display latency

		// truncate error messages to 50 chars
		// last 10 errors:
		// - error 1
		// - error 2
		// - error etc

	// Setup counter (up to 5 times) since we wait 200 milliseconds,
	// when we get to 5 it means it's time to look at the totalrequest count
	// and subtract from the previous value

	previousTotalRequestCount := 0
	cyclesCount := 0
	var rps int = 0

	for {
		sm.mu.Lock()

		// if len(sm.errors) != 0
		// append string at the end with the errors formatted
		// Only show last 10 errors
		// Errors (last 10):
		// - one error per line
		// - two errors, etc
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

		// errorsLines := "- error 1\n- error 2\n- error 3\n"

		fmt.Fprintf(sm.writer, fmtStr, sm.State, sm.linksCount, sm.depth, sm.maxDepthLevel,
			sm.linksInQueue, sm.workersRunning, sm.totalWorkersCount, sm.totalRequestsCount,
			sm.errorCounts, errorsPerc, rps, sm.lMin, sm.lAvgSum/sm.lAvgCount, sm.lMax) //, errorsLines)
		sm.mu.Unlock()

		// Only stop when AppState == Finished!
		if sm.State == AppState_Finished {
			break
		}

		time.Sleep(time.Millisecond * 200)
		cyclesCount++
	}

	// fmt.Fprintf(sm.writer, "Finished!\nTotal Links found: %d", sm.linksCount)
	sm.writer.Stop() // flush and stop rendering
}
