package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

type StatsManager struct {
	sync.Mutex

	// keep a reference to were to print stats
	writer *uilive.Writer

	state AppState

	// This is the total number of links still to be checked
	// This number will keep increasing as new links are found.
	linksInQueue   int
	linksCount     int
	errorCounts    int
	workersRunning int
	// current level of depth
	depth int

	// List of errors that happen during crawling
	errors []error

	// --------------
	// Read only vars
	totalWorkersCount int
	// depth level provided by user
	maxDepthLevel int
}

func NewStatsManager(totalWorkersCount int, depth int) *StatsManager {
	sm := StatsManager{state: AppState_IDLE, totalWorkersCount: totalWorkersCount, maxDepthLevel: depth}
	sm.writer = uilive.New()

	return &sm
}

// UpdateStats updates the stats.
func (sm *StatsManager) UpdateStats(updates ...func(*StatsManager)) {
	sm.Lock()
	defer sm.Unlock()

	for _, update := range updates {
		update(sm)
	}
}

func SetAppState(value AppState) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.state = value
	}
}

func SetLinksInQueue(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.linksInQueue = value
	}
}

func IncDecLinksInQueue(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.linksInQueue += value
	}
}

func SetLinksCount(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.linksCount = value
	}
}

func IncDecLinksCount(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.linksCount += value
	}
}

func SetErrorsCount(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.errorCounts = value
	}
}

func IncDecErrorsCount(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.errorCounts += value
	}
}

func SetWorkersRunning(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.workersRunning = value
	}
}

func IncDecWorkersRunning(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.workersRunning += value
	}
}

func SetDepth(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.depth = value
	}
}

func IncDecDepth(value int) func(*StatsManager) {
	return func(sm *StatsManager) {
		sm.depth += value
	}
}

// This functions writes to an io.Writer the updated stats
// Run this in a goroutine
func (sm *StatsManager) RunWriter() {
	sm.writer.Start()

	fmtStr := "App State: %13s      Workers: (%3d/%3d)\nLinks in Queue: %8d      Depth:     (%2d/%2d)\nLinks in Cache: %8d      Errors: %10d\n"

	for {
		sm.Lock()

		// if len(sm.errors) != 0
		// append string at the end with the errors formatted
		// Only show last 10 errors
		// Errors (last 10):
		// - one error per line
		// - two errors, etc

		fmt.Fprintf(sm.writer, fmtStr, sm.state, sm.workersRunning, sm.totalWorkersCount, sm.linksInQueue,
			sm.depth, sm.maxDepthLevel, sm.linksCount, sm.errorCounts)
		sm.Unlock()

		// Only stop when AppState == Finished!
		if sm.state == AppState_Finished {
			break
		}

		time.Sleep(time.Millisecond * 200)
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
