package internal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutePlanerParallelInvalidExecutions(t *testing.T) {
	var settings SeasonSettings
	err := ExecutePlanerParallel(settings, t.TempDir(), 0)
	assert.Error(t, err)
}

func TestExecutePlanerSerialInvalidExecutions(t *testing.T) {
	var settings SeasonSettings
	err := ExecutePlanerSerial(settings, t.TempDir(), -1)
	assert.Error(t, err)
}

func TestExecutePlanerParallelInvalidSettings(t *testing.T) {
	settings := SeasonSettings{Location: "Not/AZone"}
	err := ExecutePlanerParallel(settings, t.TempDir(), 1)
	assert.Error(t, err)
}

func TestExecutePlanerSerialSuccess(t *testing.T) {
	content, err := os.ReadFile("testdata/settings.json")
	assert.NoError(t, err)
	var settings SeasonSettings
	err = json.Unmarshal(content, &settings)
	assert.NoError(t, err)
	err = ExecutePlanerSerial(settings, t.TempDir(), 1)
	assert.NoError(t, err)
}
