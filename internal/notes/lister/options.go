package lister

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
	"github.com/jamesl33/zk/internal/note"
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
}

// matches - TODO
func (o Options) matches(n *note.Note) (bool, error) {
	// TODO
	if o.name != "" {
		return o.name == n.Name(), nil
	}

	fm, err := n.Frontmatter()
	if err != nil {
		return false, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if o.fixed != "" {
		return strings.Contains(fm.Title, o.fixed), nil
	}

	// TODO
	if o.glob != nil {
		return o.glob.Match(fm.Title), nil
	}

	// TODO
	if o.regex != nil {
		return o.regex.MatchString(fm.Title), nil
	}

	return true, nil
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
