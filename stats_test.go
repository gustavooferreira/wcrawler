package wcrawler_test

import (
	"bytes"
	"testing"

	"github.com/gustavooferreira/wcrawler"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsCLIOutWriterInstanciation(t *testing.T) {
	buf := &bytes.Buffer{}
	sm := wcrawler.NewStatsCLIOutWriter(buf, 10, 5)

	sm.SetAppState(wcrawler.AppState_Finished)
	sm.RunOutputFlusher()

	assert.Contains(t, buf.String(), "Crawler State:    Finished")
}
