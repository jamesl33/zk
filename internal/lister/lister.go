package lister

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"strings"

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

// Many returns an iterator containing matching notes.
func (l *Lister) Many(ctx context.Context) iter.Seq2[*note.Note, error] {
	return func(yield func(*note.Note, error) bool) {
		err := filepath.WalkDir(l.options.path, func(path string, _ os.DirEntry, err error) error {
			return l.walk(ctx, path, err, yield)
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

	// Ignore as it's not a note, or is hidden
	//
	// TODO (jamesl33): Make this more configurable.
	if hidden(path) || !strings.HasSuffix(path, ".md") || path == "GEMINI.md" {
		return nil
	}

	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note at %q: %w", path, err)
	}

	if l.options.matcher != nil && !l.options.matcher(n) {
		return nil
	}

	if !yield(n, nil) {
		return io.EOF
	}

	return nil
}
