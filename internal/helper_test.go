package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInSlice(t *testing.T) {
	tests := []struct {
		element interface{}
		slice   interface{}
		expected    bool
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