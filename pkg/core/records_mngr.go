package core

// RecordsManager keep track of which links we've visited and which ones we need to visit.
// Consolidates all data.
type RecordsManager struct {
	records map[string]Record
}

// This might update the stats on what it has found so far to be printed on the screen.
func (rm *RecordsManager) AddRecord() {

}

func (rm *RecordsManager) SaveToFile(file string) {

}

type Record struct {
	URL      string
	Hostname string
	Title    string
	depth    uint
	edges    []uint
}

// Entry is what gets sent to the channel for workers to pull data from the web
type Entry struct {
	URL   string
	depth uint
}
