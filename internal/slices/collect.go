package slices

import (
	"iter"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
)

// Collect2 the given iterator into a slice while handling potential errors.
func Collect2[S ~[]E, E any](seq iter.Seq2[E, error]) (S, error) {
	var s S

	err := iterator.ForEach2(seq, hs.Infallible(func(e E) {
		s = append(s, e)
	}))
	if err != nil {
		return nil, err // Purposefully not wrapped
	}

	return s, nil
}
