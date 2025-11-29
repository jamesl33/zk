package iterator

import (
	"iter"
)

// ForEach - TODO
func ForEach[T any](i iter.Seq[T], fn func(t T)) {
	for t := range i {
		fn(t)
	}
}

// ForEach2 - TODO
func ForEach2[T any](i iter.Seq2[T, error], fn func(t T) error) error {
	for t, err := range i {
		if err != nil {
			return err
		}

		err = fn(t)
		if err != nil {
			return err
		}
	}

	return nil
}
