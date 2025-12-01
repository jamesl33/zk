package hs

// Infallible wraps the given function with one which returns a <nil> error.
func Infallible[T any](fn func(t T)) func(t T) error {
	return func(t T) error { fn(t); return nil }
}
