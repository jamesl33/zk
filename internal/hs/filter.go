package hs

// Filter returns all the elements from the given slice which don't match the provided predicate.
func Filter[A any](as []A, fn func(A) bool) []A {
	af := make([]A, 0, len(as))

	for _, a := range as {
		if !fn(a) {
			continue
		}

		af = append(af, a)
	}

	return af
}
