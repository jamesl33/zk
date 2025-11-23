package lister

import (
	"regexp"

	"github.com/gobwas/glob"
)

// Options - TODO
type Options struct {
	// path - TODO
	path string

	// name - TODO
	name string

	// fixed - TODO
	fixed string

	// glob - TODO
	glob glob.Glob

	// regex - TODO
	regex *regexp.Regexp

	// tagged - TODO
	tagged []string

	// ntagged - TODO
	ntagged []string
}

// WithPath - TODO
func WithPath(path string) func(*Options) {
	return func(o *Options) {
		o.path = path
	}
}

// WithName - TODO
func WithName(name string) func(*Options) {
	return func(o *Options) {
		o.name = name
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

// WithTagged - TODO
func WithTagged(tags []string) func(*Options) {
	return func(o *Options) {
		o.tagged = tags
	}
}

// WithNotTagged - TODO
func WithNotTagged(tags []string) func(*Options) {
	return func(o *Options) {
		o.ntagged = tags
	}
}
