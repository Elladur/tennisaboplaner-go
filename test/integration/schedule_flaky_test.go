//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Elladur/tennisaboplaner-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestDefaultRun200Executions(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "settings.json"))
	assert.NoError(t, err)
	var settings internal.SeasonSettings
	err = json.Unmarshal(content, &settings)
	assert.NoError(t, err)

	fixture := loadScheduleFixture(t, filepath.Join("testdata", "schedules", "settings.json"))

	bestScore := 1e18
	results := make(chan float64, 200)
	for i := 0; i < 200; i++ {
		go func() {
			s, err := internal.CreateSeasonFromSettings(settings)
			if err != nil {
				results <- 1e18
				return
			}
			opt := internal.Optimizer{Season: &s}
			results <- opt.Optimize()
		}()
	}
	for i := 0; i < 200; i++ {
		score := <-results
		assert.NotEqual(t, 1e18, score)
		if score < bestScore {
			bestScore = score
		}
	}
	t.Logf("200-run best score: %v (500-run baseline: %v)", bestScore, fixture.Score)
	assert.LessOrEqual(t, bestScore, fixture.Score*1.01)
}
