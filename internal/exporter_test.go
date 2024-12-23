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
	season.exportExcel(tempDir)
	result, err := os.Stat(tempDir)
	if err != nil {
		t.Error("couldnt get file info")
	}
	assert.Greater(t, result.Size(), int64(0))
	assert.True(t, result.IsDir())
}

func TestExportCalendarFiles(t *testing.T) {
	season := setupStaticTestSeason()
	tempDir := t.TempDir()
	season.exportCalendarFiles(tempDir)
	result, err := os.Stat(tempDir)
	if err != nil {
		t.Error("couldnt get file info")
	}
	assert.Greater(t, result.Size(), int64(0))
	assert.True(t, result.IsDir())
}
