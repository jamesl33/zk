package ptr

// To returns a pointer to a copy of the given value.
func To[T any](v T) *T {
	return &v
}
