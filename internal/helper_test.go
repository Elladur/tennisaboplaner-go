package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInSlice(t *testing.T) {
	tests := []struct {
		element  interface{}
		slice    interface{}
		expected bool
	}{
		{element: 1, slice: []int{1, 2, 3}, expected: true},
		{element: 4, slice: []int{1, 2, 3}, expected: false},
		{element: "a", slice: []string{"a", "b", "c"}, expected: true},
		{element: "d", slice: []string{"a", "b", "c"}, expected: false},
		{element: 1.1, slice: []float64{1.1, 2.2, 3.3}, expected: true},
		{element: 4.4, slice: []float64{1.1, 2.2, 3.3}, expected: false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			switch v := tt.element.(type) {
			case int:
				assert.Equal(t, tt.expected, isInSlice(v, tt.slice.([]int)))
			case string:
				assert.Equal(t, tt.expected, isInSlice(v, tt.slice.([]string)))
			case float64:
				assert.Equal(t, tt.expected, isInSlice(v, tt.slice.([]float64)))
			}
		})
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		slice    interface{}
		expected interface{}
	}{
		{slice: []int{1, 2, 3}, expected: 3},
		{slice: []string{"a", "b", "c"}, expected: "c"},
		{slice: []float64{1.1, 2.2, 3.3}, expected: 3.3},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			switch v := tt.slice.(type) {
			case []int:
				newSlice, el := pop(v)
				assert.Equal(t, tt.expected, el)
				assert.Equal(t, v[:len(v)-1], newSlice)
			case []string:
				newSlice, el := pop(v)
				assert.Equal(t, tt.expected, el)
				assert.Equal(t, v[:len(v)-1], newSlice)
			case []float64:
				newSlice, el := pop(v)
				assert.Equal(t, tt.expected, el)
				assert.Equal(t, v[:len(v)-1], newSlice)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	tests := []struct {
		slice interface{}
	}{
		{slice: []int{1, 2, 3, 4, 5}},
		{slice: []string{"a", "b", "c", "d", "e"}},
		{slice: []float64{1.1, 2.2, 3.3, 4.4, 5.5}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			switch v := tt.slice.(type) {
			case []int:
				original := make([]int, len(v))
				copy(original, v)
				shuffled := shuffle(v)
				assert.ElementsMatch(t, v, shuffled)
				assert.NotEqual(t, shuffled, original)
			case []string:
				original := make([]string, len(v))
				copy(original, v)
				shuffled := shuffle(v)
				assert.ElementsMatch(t, v, shuffled)
				assert.NotEqual(t, shuffled, original)
			case []float64:
				original := make([]float64, len(v))
				copy(original, v)
				shuffled := shuffle(v)
				assert.ElementsMatch(t, v, shuffled)
				assert.NotEqual(t, shuffled, original)
			}
		})
	}
}
