package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Elladur/tennisaboplaner-go/internal"
)

// TestOptimalParallelRuns benchmarks how many parallel planning runs are worth
// doing: it reports the best/average/worst score and wall time for several
// candidate counts. Run with:
//
//	go test -run TestOptimalParallelRuns -v ./test/integration/...
//
// Re-run this after changing the optimizer/scoring logic to pick a new default
// for the --times flag.
func TestOptimalParallelRuns(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping run-count benchmark in short mode")
	}
	content, err := os.ReadFile("input/performance.json")
	if err != nil {
		t.Fatal(err)
	}
	var settings internal.SeasonSettings
	if err := json.Unmarshal(content, &settings); err != nil {
		t.Fatal(err)
	}

	candidates := []int{10, 50, 100, 200, 500}
	trials := 10
	t.Logf("%-6s %-8s %-10s %-10s %-10s %s", "runs", "trials", "best", "avg", "worst", "avgWall")
	for _, n := range candidates {
		scores := make([]float64, 0, trials)
		var total time.Duration
		for tr := 0; tr < trials; tr++ {
			start := time.Now()
			best := 1e18
			var mu sync.Mutex
			var wg sync.WaitGroup
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s, err := internal.CreateSeasonFromSettings(settings)
					if err != nil {
						t.Error(err)
						return
					}
					score := (&internal.Optimizer{Season: &s}).Optimize()
					mu.Lock()
					if score < best {
						best = score
					}
					mu.Unlock()
				}()
			}
			wg.Wait()
			scores = append(scores, best)
			total += time.Since(start)
		}
		min, max, sum := scores[0], scores[0], 0.0
		for _, v := range scores {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
			sum += v
		}
		t.Logf("%-6d %-8d %-10.6f %-10.6f %-10.6f %s", n, trials, min, sum/float64(trials), max, total/time.Duration(trials))
	}
	_ = fmt.Sprintf
}
