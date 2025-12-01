package iterator

import (
	"iter"
)

// Empty returns an empty iterator.
func Empty[T any]() iter.Seq2[T, error] {
	return func(_ func(T, error) bool) {}
}
