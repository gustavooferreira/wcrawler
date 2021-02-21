package core

// RecordManager keep track of which links we've visited and which ones we need to visit.
// Consolidates all data.
type RecordManager struct {
	// Keeps a table of records. Key is the URL (scheme,authority,path,query)
	records    map[string]Record
	indexCount uint
}

func NewRecordManager() *RecordManager {
	records := make(map[string]Record)
	rm := RecordManager{records: records}
	return &rm
}

// This might update the stats on what it has found so far to be printed on the screen.
// No, the merger goroutine will update the stats
func (rm *RecordManager) AddRecord(entry Entry) {
	// Check for each URL if they already exist
	// if not add them to the table including a new index
	// reference that index in the ParentURL edges slice

}

// Match checks whether this URL already exists in the DB
func (rm *RecordManager) Match(rawURL string) bool {
	_, ok := rm.records[rawURL]
	return ok
}

// SaveToFile dumps the records map into a file in JSON format.
func (rm *RecordManager) SaveToFile(filepath string) {

}

// LoadFromFile dumps the records map into a file in JSON format.
func (rm *RecordManager) LoadFromFile(filepath string) {

}
