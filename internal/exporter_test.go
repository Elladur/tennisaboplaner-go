package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestExport(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	err := season.Export(tempDir)
	assert.NoError(t, err)
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
	err := season.exportCalendarFiles(tempDir)
	assert.NoError(t, err)
	for _, player := range season.Players {
		result, err := os.Stat(tempDir + string(os.PathSeparator) + player.Name + ".ics")
		if err != nil {
			t.Error("couldnt get file info for player " + player.Name)
		}
		assert.Greater(t, result.Size(), int64(0))
		assert.False(t, result.IsDir())
	}
}

func TestExportExcelHasGapStatsSheet(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	err := season.exportExcel(tempDir)
	assert.NoError(t, err)
	f, err := excelize.OpenFile(tempDir + string(os.PathSeparator) + "schedule.xlsx")
	assert.NoError(t, err)
	defer func() { _ = f.Close() }()
	idx, err := f.GetSheetIndex(gapStatsSheet)
	assert.NoError(t, err)
	assert.NotEqual(t, -1, idx, "missing sheet %s", gapStatsSheet)
}

func TestGapLengths(t *testing.T) {
	assert.Equal(t, []int{0, 0, 0, 0, 0, 0}, gapLengths(5, []int{0, 1, 2, 3, 4}))
	assert.Equal(t, []int{1, 1, 1}, gapLengths(5, []int{1, 3}))
	assert.Equal(t, []int{0, 3, 0}, gapLengths(5, []int{0, 4}))
	assert.Equal(t, []int{5}, gapLengths(5, []int{}))
	assert.Equal(t, []int{0, 4}, gapLengths(5, []int{0}))
	assert.Equal(t, []int{2, 2, 2}, gapLengths(8, []int{2, 5}))
}
