package internal

import (
	"math/rand"
	"path/filepath"
	"strings"
)

// sanitizeFilename replaces characters that are unsafe in file names and
// falls back to a placeholder for empty results.
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "player"
	}
	repl := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_",
		"|", "_", "\x00", "_",
	)
	name = repl.Replace(name)
	name = strings.Trim(name, " .")
	if name == "" {
		name = "player"
	}
	// Guard against path traversal after cleaning.
	name = filepath.Base(name)
	if name == "." || name == ".." || name == "" {
		name = "player"
	}
	return name
}

func isInSlice[T comparable](element T, slice []T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func pop[T any](slice []T) ([]T, T) {
	el := slice[len(slice)-1]
	slice = slice[:len(slice)-1]
	return slice, el
}

// shuffle mutates and returns the given slice.
func shuffle[T any](slice []T) []T {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
