package lister

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"github.com/jamesl33/zk/internal/note"
)

// Lister is a note lister which iterates directories recursively finding matching notes.
type Lister struct {
	options Options
}

// NewLister returns an initialized lister.
func NewLister(opts ...func(o *Options)) (*Lister, error) {
	var o Options

	for _, opt := range opts {
		opt(&o)
	}

	lister := Lister{
		options: o,
	}

	return &lister, nil
}

// One returns the first match for the lister.
func (l *Lister) One(ctx context.Context) (*note.Note, error) {
	next, stop := iter.Pull2(l.Many(ctx))
	defer stop()

	n, err, ok := next()
	if !ok {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get next note: %w", err)
	}

	return n, nil
}

// Many returns an iterator containing matching notes.
func (l *Lister) Many(ctx context.Context) iter.Seq2[*note.Note, error] {
	return func(yield func(*note.Note, error) bool) {
		err := filepath.WalkDir(l.options.path, func(path string, entry os.DirEntry, err error) error {
			return l.walk(ctx, path, entry, err, yield)
		})
		if err == nil || errors.Is(err, io.EOF) {
			return
		}

		yield(nil, err)
	}
}

// walk the given directory finding matching notes.
func (l *Lister) walk(
	ctx context.Context,
	path string,
	entry os.DirEntry,
	err error,
	yield func(n *note.Note, err error) bool,
) error {
	// Exit early as the walk has been canceled.
	if err := ctx.Err(); err != nil {
		return err // Purposefully not wrapped
	}

	if err != nil {
		return fmt.Errorf("unexpected error walking %q: %w", path, err)
	}

	var (
		hidden = hidden(path)
		ignore = ignore(path)
	)

	// Ignore the directory; it's hidden
	if entry.IsDir() && hidden {
		return filepath.SkipDir
	}

	// Ignore the note as it's hidden, or ignored
	if hidden || ignore {
		return nil
	}

	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note at %q: %w", path, err)
	}

	m, err := l.match(n)
	if err != nil {
		return err // Purposefully not wrapped
	}

	if !m {
		return nil
	}

	if !yield(n, nil) {
		return io.EOF
	}

	return nil
}

// match returns a boolean indicating whether the given note should be listed.
func (l *Lister) match(n *note.Note) (bool, error) {
	if l.options.matcher == nil {
		return true, nil
	}

	m, err := l.options.matcher(n)
	if err != nil {
		return false, fmt.Errorf("failed to run matcher: %w", err)
	}

	return m, nil
}
