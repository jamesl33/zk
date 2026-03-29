package hs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDifference(t *testing.T) {
	type test struct {
		name     string
		as       []string
		bs       []string
		expected []string
	}

	tests := []test{
		{
			name:     "no overlap",
			as:       []string{"a", "b"},
			bs:       []string{"c", "d"},
			expected: []string{"a", "b"},
		},
		{
			name:     "partial overlap",
			as:       []string{"a", "b", "c"},
			bs:       []string{"b", "c", "d"},
			expected: []string{"a"},
		},
		{
			name:     "full overlap",
			as:       []string{"a", "b"},
			bs:       []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "empty as",
			as:       []string{},
			bs:       []string{"a", "b"},
			expected: []string{},
		},
		{
			name:     "empty bs",
			as:       []string{"a", "b"},
			bs:       []string{},
			expected: []string{"a", "b"},
		},
		{
			name:     "both empty",
			as:       []string{},
			bs:       []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Difference(tt.as, tt.bs))
		})
	}
}
