package iterator

import (
	"testing"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/stretchr/testify/assert"
)

func TestForEach(t *testing.T) {
	type test struct {
		name  string
		input []string
	}

	tests := []test{
		{
			name:  "MultipleItems",
			input: []string{"a", "b", "c"},
		},
		{
			name:  "EmptyItems",
			input: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := func(yield func(string) bool) {
				for _, item := range tt.input {
					if !yield(item) {
						return
					}
				}
			}

			actual := make([]string, 0)

			ForEach(it, func(s string) {
				actual = append(actual, s)
			})

			assert.Equal(t, tt.input, actual)
		})
	}
}

func TestForEach2(t *testing.T) {
	type test struct {
		name  string
		input []string
	}

	tests := []test{
		{
			name:  "MultipleItems",
			input: []string{"a", "b", "c"},
		},
		{
			name:  "EmptyItems",
			input: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := func(yield func(string, error) bool) {
				for _, item := range tt.input {
					if !yield(item, nil) {
						return
					}
				}
			}

			actual := make([]string, 0)

			err := ForEach2(it, hs.Infallible(func(s string) {
				actual = append(actual, s)
			}))

			assert.NoError(t, err)
			assert.Equal(t, tt.input, actual)
		})
	}
}

func TestForEach2IteratorError(t *testing.T) {
	it := func(yield func(string, error) bool) {
		yield("", assert.AnError)
	}

	err := ForEach2(it, func(s string) error {
		return nil
	})

	assert.ErrorIs(t, err, assert.AnError)
}

func TestForEach2FnError(t *testing.T) {
	it := func(yield func(string, error) bool) {
		yield("a", nil)
	}

	err := ForEach2(it, func(s string) error {
		return assert.AnError
	})

	assert.ErrorIs(t, err, assert.AnError)
}
