package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	season.Export(tempDir)
	result, err := os.Stat(tempDir)
	if err != nil {
		t.Error("couldnt get file info")
	}
	assert.Greater(t, result.Size(), int64(0))
	assert.True(t, result.IsDir())
}

func TestExportExcel(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	err := season.exportExcel(tempDir)
	assert.NoError(t, err)
	result, err := os.Stat(tempDir + string(os.PathSeparator) + "schedule.xlsx")
	if err != nil {
		t.Error("couldnt get file info")
	}
	assert.Greater(t, result.Size(), int64(0))
	assert.False(t, result.IsDir())
}

func TestExportCalendarFiles(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	season.exportCalendarFiles(tempDir)
	for _, player := range season.Players {
		result, err := os.Stat(tempDir + string(os.PathSeparator) + player.Name + ".ics")
		if err != nil {
			t.Error("couldnt get file info for player " + player.Name)
		}
		assert.Greater(t, result.Size(), int64(0))
		assert.False(t, result.IsDir())
	}
}
