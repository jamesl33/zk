package iterator

import (
	"fmt"
	"iter"
)

// ForEach runs the given function for each element in the provided iterator.
func ForEach[T any](i iter.Seq[T], fn func(t T)) {
	for t := range i {
		fn(t)
	}
}

// ForEach2 runs the given function for each element in the provided iterator.
func ForEach2[T any](i iter.Seq2[T, error], fn func(t T) error) error {
	for t, err := range i {
		if err != nil {
			return fmt.Errorf("unexpected error during iteration: %w", err)
		}

		err = fn(t)
		if err != nil {
			return err // Purposefully not wrapped
		}
	}

	return nil
}
