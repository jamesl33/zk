package searcher

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

// Searcher - TODO
type Searcher struct {
	options Options
}

// NewSearcher - TODO
func NewSearcher(opts ...func(o *Options)) (*Searcher, error) {
	var o Options

	for _, opt := range opts {
		opt(&o)
	}

	lister := Searcher{
		options: o,
	}

	return &lister, nil
}

// One - TODO
func (s *Searcher) One(ctx context.Context) (*note.Note, error) {
	next, stop := iter.Pull2(s.Many(ctx))
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
func (s *Searcher) Many(ctx context.Context) iter.Seq2[*note.Note, error] {
	return func(yield func(*note.Note, error) bool) {
		err := filepath.WalkDir(s.options.path, func(path string, _ os.DirEntry, err error) error {
			return s.walk(ctx, path, err, yield)
		})
		if err == nil {
			return
		}

		yield(nil, err)
	}
}

// walk - TODO
func (s *Searcher) walk(
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

	match, err := s.matches(n)
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
func (s *Searcher) matches(n *note.Note) (bool, error) {
	body, err := n.Body()
	if err != nil {
		return false, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if s.options.fixed != "" {
		return strings.Contains(body, s.options.fixed), nil
	}

	// TODO
	if s.options.glob != nil {
		return s.options.glob.Match(body), nil
	}

	// TODO
	if s.options.regex != nil {
		return s.options.regex.MatchString(body), nil
	}

	return true, nil
}
