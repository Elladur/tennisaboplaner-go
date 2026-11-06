package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUpdateAvailable(t *testing.T) {
	cases := []struct {
		current, latest string
		want            bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "1.0.0", false},
		{"1.2.0", "1.1.9", false},
		{"dev", "1.0.0", true},
		{"1.0.0", "not-a-version", false},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, isUpdateAvailable(c.current, c.latest), "%s -> %s", c.current, c.latest)
	}
}

func TestVersionDefault(t *testing.T) {
	assert.Equal(t, "dev", Version)
}
