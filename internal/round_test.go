package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPlayersOfRound(t *testing.T) {
	tests := []struct {
		name     string
		round    []Match
		expected []int
	}{
		{
			name:     "No matches",
			round:    []Match{},
			expected: []int{},
		},
		{
			name: "Single match",
			round: []Match{
				{1, 2, true},
			},
			expected: []int{1, 2},
		},
		{
			name: "Multiple matches",
			round: []Match{
				{1, 2, true},
				{3, 4, true},
			},
			expected: []int{1, 2, 3, 4},
		},
		{
			name: "Overlapping players",
			round: []Match{
				{1, 2, true},
				{2, 3, true},
			},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPlayersOfRound(tt.round)
			assert.Equal(t, tt.expected, result)
		})
	}
}
