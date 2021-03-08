package stats_test

import (
	"testing"

	"github.com/gustavooferreira/wcrawler/internal/stats"
	"github.com/stretchr/testify/assert"
)

func TestNewStatsManagerCreation(t *testing.T) {
	sm := stats.NewStatsManager(10, 5)

	sm.UpdateStats(stats.SetAppState(stats.AppState_Running))

	assert.Equal(t, stats.AppState_Running, int(sm.State))
}
