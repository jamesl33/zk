package hs

// Infallible - TODO
func Infallible[T any](fn func(t T)) func(t T) error {
	return func(t T) error { fn(t); return nil }
}
