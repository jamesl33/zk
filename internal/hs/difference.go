package hs

// Difference returns the elements in `as` that are not in `bs`.
func Difference[T comparable](as, bs []T) []T {
	lu := make(map[T]struct{})

	for _, b := range bs {
		lu[b] = struct{}{}
	}

	return Filter(as, func(a T) bool { _, ok := lu[a]; return !ok })
}
