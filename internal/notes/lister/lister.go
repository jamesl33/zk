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

// Lister - TODO
type Lister struct {
	options Options
}

// NewLister - TODO
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

// Many - TODO
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

// walk - TODO
func (l *Lister) walk(
	ctx context.Context,
	path string,
	err error,
	yield func(n *note.Note, err error) bool,
) error {
	// TODO
	if err := ctx.Err(); err != nil {
		return err // Purposefully not wrapped
	}

	// TODO
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if hidden(path) || !strings.HasSuffix(path, ".md") {
		return nil
	}

	// TODO
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	if l.options.matcher != nil && !l.options.matcher(n) {
		return nil
	}

	if !yield(n, nil) {
		return io.EOF
	}

	return nil
}
