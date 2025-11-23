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

// yielder - TODO
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

// One - TODO
func (l *Lister) One(ctx context.Context) (*note.Note, error) {
	next, stop := iter.Pull2(l.Many(ctx))
	defer stop()

	n, err, ok := next()
	if !ok {
		return nil, errors.New("not found") // TODO
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return n, nil
}

// Many - TODO
func (l *Lister) Many(ctx context.Context) iter.Seq2[*note.Note, error] {
	return func(yield func(*note.Note, error) bool) {
		err := filepath.WalkDir(l.options.path, func(path string, _ os.DirEntry, err error) error {
			return l.walk(ctx, path, err, yield)
		})
		if err == nil {
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
	n := note.NewNote(path)

	match, err := l.matches(n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	if !match {
		return nil
	}

	if !yield(n, nil) {
		return io.EOF
	}

	return nil
}

// matches - TODO
func (l *Lister) matches(n *note.Note) (bool, error) {
	// TODO
	if l.options.name != "" {
		return l.options.name == n.Name(), nil
	}

	fm, err := n.Frontmatter()
	if err != nil {
		return false, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if l.options.fixed != "" {
		return strings.Contains(fm.Title, l.options.fixed), nil
	}

	// TODO
	if l.options.glob != nil {
		return l.options.glob.Match(fm.Title), nil
	}

	// TODO
	if l.options.regex != nil {
		return l.options.regex.MatchString(fm.Title), nil
	}

	return true, nil
}
