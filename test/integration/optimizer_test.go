package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Elladur/tennisaboplaner-go/internal"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestPerformance(t *testing.T) {
	log.SetLevel(log.WarnLevel) // to make it faster and ignore logs
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

	start := time.Now()
	internal.ExecutePlanerSerial(settings, t.TempDir(), 100)
	elapsed := time.Since(start)
	fmt.Printf("100 optimizations took %s\n", elapsed)
	assert.Less(t, elapsed, 6*time.Second)
}
