package hs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	type test struct {
		name     string
		input    []int
		fn       func(int) bool
		expected []int
	}

	tests := []test{
		{
			name:     "filter evens",
			input:    []int{1, 2, 3, 4, 5, 6},
			fn:       func(a int) bool { return a%2 == 0 },
			expected: []int{2, 4, 6},
		},
		{
			name:     "filter none",
			input:    []int{1, 2, 3},
			fn:       func(a int) bool { return true },
			expected: []int{1, 2, 3},
		},
		{
			name:     "filter all",
			input:    []int{1, 2, 3},
			fn:       func(a int) bool { return false },
			expected: []int{},
		},
		{
			name:     "empty input",
			input:    []int{},
			fn:       func(a int) bool { return true },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Filter(tt.input, tt.fn))
		})
	}
}
