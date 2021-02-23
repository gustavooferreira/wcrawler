package core

import (
	"encoding/json"
	"io"
)

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
func (rm *RecordManager) Get(rawURL string) (Record, bool) {
	if r, ok := rm.records[rawURL]; ok {
		return r, true
	}
	return Record{}, false
}

// Count counts the number of records.
func (rm *RecordManager) Count(rawURL string) int {
	return len(rm.records)
}

// Dump returns all records in the RecordManager.
func (rm *RecordManager) Dump(rawURL string) map[string]Record {
	return rm.records
}

// SaveToWriter dumps the records map into a Writer in JSON format.
// Can pass a os.File, to write to a file.
func (rm *RecordManager) SaveToWriter(w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(rm.records)
	return err
}

// LoadFromReader reads the records from a Reader in JSON format.
// Can pass a os.File, to read from a file.
func (rm *RecordManager) LoadFromReader(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&rm.records)
	return err
}
