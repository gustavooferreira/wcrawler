package wcrawler

import "fmt"

// AppState represents the current state of the App.
type AppState int

const (
	// AppState_Unknown represents the 'unknown' state.
	AppState_Unknown = iota
	// AppState_IDLE represents the 'idle' state.
	AppState_IDLE
	// AppState_Running represents the 'run' state.
	AppState_Running
	// AppState_Finished represents the 'finish' state.
	AppState_Finished
)

var appStateToString = map[AppState]string{
	AppState_Unknown:  "Unknown",
	AppState_IDLE:     "IDLE",
	AppState_Running:  "Running",
	AppState_Finished: "Finished",
}

var appStateToEnum = map[string]AppState{
	"Unknown":  AppState_Unknown,
	"IDLE":     AppState_IDLE,
	"Running":  AppState_Running,
	"Finished": AppState_Finished,
}

// String returns the string representation of AppState.
func (as AppState) String() string {
	state, ok := appStateToString[as]
	if !ok {
		return "Unknown"
	}

	return state
}

// Parse parses a string into AppState returning an error if string passed cannot be parsed into a valid state.
func (as *AppState) Parse(state string) error {
	value, ok := appStateToEnum[state]
	if !ok {
		return fmt.Errorf("couldn't parse app state")
	}

	*as = value
	return nil
}
