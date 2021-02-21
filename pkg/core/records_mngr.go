package core

// RecordsManager keep track of which links we've visited and which ones we need to visit.
// Consolidates all data.
type RecordsManager struct {
	records map[string]Record
}

// This might update the stats on what it has found so far to be printed on the screen.
func (rm *RecordsManager) AddRecord() {

}

// Match checks whether this URL already exists in the DB
func (rm *RecordsManager) Match() bool {
	return false
}

func (rm *RecordsManager) SaveToFile(file string) {

}
