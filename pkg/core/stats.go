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
	linksInQueue   uint
	linksCount     uint
	errorCounts    uint
	workersRunning uint
	// current level of depth
	depth uint

	// List of errors that happen during crawling
	errors []error

	// --------------
	// Read only vars
	totalWorkersCount uint
	// depth level provided by user
	maxDepthLevel uint
}

func NewStatsManager(totalWorkersCount uint, depth uint) *StatsManager {
	sm := StatsManager{state: AppState_IDLE, totalWorkersCount: totalWorkersCount, maxDepthLevel: depth}
	sm.writer = uilive.New()

	return &sm
}

// This updates the stats counters, this is cumulative, meaning only put the numbers to add to the total
// workersRunningCounter can be negative, when they finish processing they will decrement this.
func (sm *StatsManager) UpdateStats(state AppState, linksInQueueInc uint, linksCountInc uint, errorCountsInc uint, workersRunningCounter int) {
	sm.Lock()
	defer sm.Unlock()

	// if AppState == 0 then no update on that
	sm.workersRunning += uint(workersRunningCounter)

}

// This functions writes to an io.Writer the updated stats
// Run this in a goroutine
func (sm *StatsManager) RunWriter() {
	sm.writer.Start()

	for {
		sm.Lock()
		fmt.Fprintf(sm.writer, "App State: %s      Workers (%d/%d)\nLinks Count: %5d   Errors: (%d/%d)\n",
			sm.state, sm.workersRunning, sm.totalWorkersCount, sm.linksCount,
			sm.errorCounts, sm.linksCount)
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
