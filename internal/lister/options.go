package lister

import (
	"github.com/jamesl33/zk/internal/matcher"
)

// Options - TODO
type Options struct {
	// path - TODO
	path string

	// matcher - TODO
	matcher matcher.Matcher
}

// WithPath - TODO
func WithPath(path string) func(*Options) {
	return func(o *Options) {
		o.path = path
	}
}

// WithMatcher - TODO
func WithMatcher(matcher matcher.Matcher) func(*Options) {
	return func(o *Options) {
		o.matcher = matcher
	}
}
