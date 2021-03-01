package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type StatsCLIOutput struct {
	sync.Mutex

	// keep a reference to were to print stats
	writer *uilive.Writer

	state AppState

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
	lAvgTotal float64
	lAvgCount float64

	// Requests per second
	rps int

	// List of errors that happen during crawling
	errors [10]error

	// --------------
	// Read only vars
	totalWorkersCount int
	// depth level provided by user
	maxDepthLevel int
}

func NewStatsManager(totalWorkersCount int, depth int) *StatsCLIOutput {
	sm := StatsCLIOutput{state: AppState_IDLE, totalWorkersCount: totalWorkersCount, maxDepthLevel: depth}
	sm.writer = uilive.New()

	return &sm
}

// UpdateStats updates the stats.
func (sm *StatsCLIOutput) UpdateStats(updates ...func(*StatsCLIOutput)) {
	sm.Lock()
	defer sm.Unlock()

	for _, update := range updates {
		update(sm)
	}
}

func SetAppState(value AppState) func(*StatsCLIOutput) {
	return func(sm *StatsCLIOutput) {
		sm.state = value
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
			sm.lAvgTotal = valueF
		} else {
			if valueF > sm.lMax {
				sm.lMax = valueF
			}

			if valueF < sm.lMin {
				sm.lMin = valueF
			}

			sm.lAvgCount++
			sm.lAvgTotal += valueF
		}
	}
}

// This functions writes to an io.Writer the updated stats
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
		sm.Lock()

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

		fmt.Fprintf(sm.writer, fmtStr, sm.state, sm.linksCount, sm.depth, sm.maxDepthLevel,
			sm.linksInQueue, sm.workersRunning, sm.totalWorkersCount, sm.totalRequestsCount,
			sm.errorCounts, errorsPerc, rps, sm.lMin, sm.lAvgTotal/sm.lAvgCount, sm.lMax) //, errorsLines)
		sm.Unlock()

		// Only stop when AppState == Finished!
		if sm.state == AppState_Finished {
			break
		}

		time.Sleep(time.Millisecond * 200)
		cyclesCount++
	}

	// fmt.Fprintf(sm.writer, "Finished!\nTotal Links found: %d", sm.linksCount)
	sm.writer.Stop() // flush and stop rendering
}

// AppState represents the current state of the App.
type AppState uint

const (
	// AppState_IDLE represents the 'idle' state.
	AppState_IDLE = iota + 1
	// AppState_Running represents the 'run' state.
	AppState_Running
	// AppState_Finished represents the 'finish' state.
	AppState_Finished
)

// String returns the string representation of AppState.
func (as AppState) String() string {
	return [...]string{"", "IDLE", "Running", "Finished"}[as]
}
