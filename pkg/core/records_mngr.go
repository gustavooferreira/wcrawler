package core

import "fmt"

// RecordManager keep track of which links we've visited and which ones we need to visit.
// Consolidates all data.
type RecordManager struct {
	// Keeps a table of records. Key is the URL (scheme,authority,path,query)
	records    map[string]Record
	indexCount uint
}

// NewRecordManager returns a new Record Manager.
func NewRecordManager() *RecordManager {
	records := make(map[string]Record)
	rm := RecordManager{records: records}
	return &rm
}

// AddRecord adds a record to the RecordManager.
func (rm *RecordManager) AddRecord(entry RMEntry) {
	var index uint

	if entryRecord, ok := rm.records[entry.URL.String]; ok {
		index = entryRecord.Index
	} else {
		index = rm.indexCount

		r := Record{
			Index:     rm.indexCount,
			ParentURL: entry.ParentURL,
			URL:       entry.URL.String,
			Host:      entry.URL.Host,
			Depth:     entry.Depth,
			Edges:     []uint{},
		}

		rm.records[entry.URL.String] = r
	}

	rm.indexCount++

	// Add pointers on parent's entry
	if entry.ParentURL != "" {
		if parentEntry, ok := rm.records[entry.ParentURL]; ok {
			parentEntry.Edges = append(parentEntry.Edges, index)
			rm.records[entry.ParentURL] = parentEntry
		} else {
			// we should have never landed here. being here, means there is a bug somewhere else.
		}
	}
}

// Match checks whether this URL already exists in the DB
func (rm *RecordManager) Match(rawURL string) bool {
	_, ok := rm.records[rawURL]
	return ok
}

// Get returns a record from the Record Manager.
func (rm *RecordManager) Get(rawURL string) (Record, error) {
	if r, ok := rm.records[rawURL]; ok {
		return r, nil
	}
	return Record{}, fmt.Errorf("no record found")
}

// Count counts the number of records.
func (rm *RecordManager) Count(rawURL string) int {
	return len(rm.records)
}

// SaveToFile dumps the records map into a file in JSON format.
func (rm *RecordManager) SaveToFile(filepath string) {

}

// LoadFromFile dumps the records map into a file in JSON format.
func (rm *RecordManager) LoadFromFile(filepath string) {

}
