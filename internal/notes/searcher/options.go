package searcher

import (
	"regexp"

	"github.com/gobwas/glob"
)

// Options - TODO
type Options struct {
	// path - TODO
	path string

	// fixed - TODO
	fixed string

	// glob - TODO
	glob glob.Glob

	// regex - TODO
	regex *regexp.Regexp
}

// WithPath - TODO
func WithPath(path string) func(*Options) {
	return func(o *Options) {
		o.path = path
	}
}

// WithFixed - TODO
func WithFixed(fixed string) func(*Options) {
	return func(o *Options) {
		o.fixed = fixed
	}
}

// WithGlob - TODO
func WithGlob(glob glob.Glob) func(*Options) {
	return func(o *Options) {
		o.glob = glob
	}
}

// WithRegex - TODO
func WithRegex(regex *regexp.Regexp) func(*Options) {
	return func(o *Options) {
		o.regex = regex
	}
}
