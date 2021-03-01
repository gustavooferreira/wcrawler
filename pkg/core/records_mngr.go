package core

import (
	"encoding/json"
	"fmt"
	"io"
)

// RecordManager keeps track of links visited and some metadata like depth level and its children.
type RecordManager struct {
	// Keeps a table of Records. Key is the URL (scheme,authority,path,query)
	Records    map[string]Record
	IndexCount uint
}

// NewRecordManager returns a new Record Manager.
func NewRecordManager() *RecordManager {
	records := make(map[string]Record)
	rm := RecordManager{Records: records}
	return &rm
}

// AddRecord adds a record to the RecordManager.
func (rm *RecordManager) AddRecord(entry RMEntry) {
	var index uint

	if entryRecord, ok := rm.Records[entry.URL.Raw]; ok {
		index = entryRecord.Index
	} else {
		index = rm.IndexCount

		r := Record{
			Index:      index,
			ParentURL:  entry.ParentURL,
			URL:        entry.URL.Raw,
			Host:       entry.URL.Host,
			Depth:      entry.Depth,
			StatusCode: entry.StatusCode,
			ErrString:  entry.ErrString,
			Edges:      []uint{},
		}

		rm.Records[entry.URL.Raw] = r
	}

	rm.IndexCount++

	// Add pointers on parent's entry
	if entry.ParentURL != "" {
		if parentEntry, ok := rm.Records[entry.ParentURL]; ok {
			parentEntry.Edges = append(parentEntry.Edges, index)
			rm.Records[entry.ParentURL] = parentEntry
		} else {
			// we should have never landed here. being here, means there is a bug somewhere else.
		}
	}
}

// Exists checks whether this URL exists in the table.
func (rm *RecordManager) Exists(rawURL string) bool {
	_, ok := rm.Records[rawURL]
	return ok
}

// Update updates entry in the table.
func (rm *RecordManager) Update(rawURL string, statusCode int, err error) error {
	if elem, ok := rm.Records[rawURL]; ok {
		elem.StatusCode = statusCode

		if err != nil {
			elem.ErrString = err.Error()
		}

		rm.Records[rawURL] = elem
		return nil
	}
	return fmt.Errorf("record not found")
}

// Get returns a record from the Record Manager.
func (rm *RecordManager) Get(rawURL string) (Record, bool) {
	r, ok := rm.Records[rawURL]
	return r, ok
}

// Count counts the number of records.
func (rm *RecordManager) Count() int {
	return len(rm.Records)
}

// Dump returns all records in the RecordManager.
func (rm *RecordManager) Dump() map[string]Record {
	return rm.Records
}

// SaveToWriter dumps the records map into a Writer in JSON format.
// Can pass a os.File, to write to a file.
func (rm *RecordManager) SaveToWriter(w io.Writer, indent bool) error {
	encoder := json.NewEncoder(w)
	if indent {
		encoder.SetIndent("", "    ")
	}
	err := encoder.Encode(rm.Records)
	return err
}

// LoadFromReader reads the records from a Reader in JSON format.
// Can pass a os.File, to read from a file.
func (rm *RecordManager) LoadFromReader(r io.Reader) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&rm.Records)
	rm.IndexCount = uint(len(rm.Records))
	return err
}
