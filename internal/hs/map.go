package hs

// Map - TODO
func Map[A any, B any](as []A, fn func(A) B) []B {
	bs := make([]B, 0, len(as))

	for _, a := range as {
		bs = append(bs, fn(a))
	}

	return bs
}
