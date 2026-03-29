package hs

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type test struct {
		name     string
		input    []int
		fn       func(int) string
		expected []string
	}

	tests := []test{
		{
			name:     "map integers to strings",
			input:    []int{1, 2, 3},
			fn:       func(a int) string { return strconv.Itoa(a) },
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "empty input",
			input:    []int{},
			fn:       func(a int) string { return strconv.Itoa(a) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Map(tt.input, tt.fn))
		})
	}
}
