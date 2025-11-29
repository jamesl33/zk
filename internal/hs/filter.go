package hs

// Filter - TODO
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
