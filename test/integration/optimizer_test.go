package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Elladur/tennisaboplaner-go/internal"
	"github.com/stretchr/testify/assert"
)

func TestPerformance(t *testing.T) {
	content, err := os.ReadFile("input/performance.json")
	if err != nil {
		t.Error("Could not read settings.json")
		return
	}
	var settings internal.SeasonSettings
	err = json.Unmarshal(content, &settings)
	if err != nil {
		t.Error("Could not unmarshal settings.json")
		return
	}
	season, err := internal.CreateSeasonFromSettings(settings)
	if err != nil {
		t.Error("Could not create season from settings")
		return
	}

	optimizer := internal.Optimizer{Season: &season}

	start := time.Now()
	for i := 0; i < 1000; i++ {
		optimizer.Optimize()
		season.CreateSchedule()
	}
	elapsed := time.Since(start)
	fmt.Printf("100 optimizations took %s\n", elapsed)
	assert.Less(t, elapsed, 11*time.Second)
}
