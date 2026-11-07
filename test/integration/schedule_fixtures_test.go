package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Elladur/tennisaboplaner-go/internal"
	"github.com/stretchr/testify/assert"
)

type scheduleFixture struct {
	Players  []internal.Player  `json:"players"`
	Schedule [][]internal.Match `json:"schedule"`
	Score    float64            `json:"score"`
}

func loadScheduleFixture(t *testing.T, path string) scheduleFixture {
	t.Helper()
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	var f scheduleFixture
	err = json.Unmarshal(content, &f)
	assert.NoError(t, err)
	return f
}

func TestFixtureScheduleScores(t *testing.T) {
	dir := filepath.Join("testdata", "schedules")
	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			f := loadScheduleFixture(t, filepath.Join(dir, e.Name()))
			assert.NotEmpty(t, f.Players)
			assert.NotEmpty(t, f.Schedule)

			got := internal.GetScore(f.Schedule, f.Players)
			assert.InDelta(t, f.Score, got, f.Score*0.05, "score should be within 5% of stored baseline")

			for _, round := range f.Schedule {
				assert.NotEmpty(t, round)
			}
		})
	}
}
