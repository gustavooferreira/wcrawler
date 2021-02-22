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

	linksCount     uint
	errorCounts    uint
	workersRunning uint
}

func NewStatsManager() *StatsManager {
	sm := StatsManager{state: AppState_IDLE}
	sm.writer = uilive.New()

	return &sm
}

// This updates the stats counters, this is cumulative, meaning only put the numbers to add to the total
// workersRunning can be negative, when they finish processing they will decrement this.
func (sm *StatsManager) UpdateStats(state AppState, linksCount uint, errorCounts uint, workersRunning int) {
	sm.Lock()
	defer sm.Unlock()

	// if AppState == 0 then no update on that

}

// This prints to an io.Reader the updated stats
// Run this in a goroutine
func (sm *StatsManager) RunWriter() {

	// start listening to updates and render
	sm.writer.Start()

	for {

		// Only stop when AppState == Finished!

		sm.Lock()
		fmt.Fprintf(sm.writer, "Downloading.. (%d/%d) GB\nDownloading.. (%d/%d) GB\n", 1, 100, 1*2, 200)
		sm.Unlock()
		time.Sleep(time.Millisecond * 50)
	}

	fmt.Fprintln(sm.writer, "Finished: Downloaded 100GB")
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
