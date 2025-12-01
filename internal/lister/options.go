package lister

import (
	"github.com/jamesl33/zk/internal/matcher"
)

// Options encapsulates the options for note listing.
type Options struct {
	// path to start listing from.
	path string

	// matcher which determines if a note should be included in the list.
	matcher matcher.Matcher
}

// WithPath allows specifying the path to start listing from.
func WithPath(path string) func(*Options) {
	return func(o *Options) {
		o.path = path
	}
}

// WithMatcher allows specifying the note matcher.
func WithMatcher(matcher matcher.Matcher) func(*Options) {
	return func(o *Options) {
		o.matcher = matcher
	}
}
