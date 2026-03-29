package slices

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect2(t *testing.T) {
	type test struct {
		name     string
		seq      iter.Seq2[int, error]
		expected []int
	}

	tests := []test{
		{
			name: "empty sequence",
			seq:  func(yield func(int, error) bool) {},
		},
		{
			name: "multiple elements",
			seq: func(yield func(int, error) bool) {
				if !yield(1, nil) {
					return
				}

				if !yield(2, nil) {
					return
				}

				if !yield(3, nil) {
					return
				}
			},
			expected: []int{1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Collect2[[]int](test.seq)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestCollect2WithError(t *testing.T) {
	seq := func(yield func(int, error) bool) {
		if !yield(1, nil) {
			return
		}

		if !yield(0, assert.AnError) {
			return
		}
	}

	actual, err := Collect2[[]int](seq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
