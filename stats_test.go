package wcrawler_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsManagerCreation(t *testing.T) {
	sm := wcrawler.NewStatsManager(10, 5)

	sm.UpdateStats(wcrawler.SetAppState(wcrawler.AppState_Running))

	assert.Equal(t, wcrawler.AppState_Running, int(sm.State))
}
