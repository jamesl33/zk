package notes

import (
	"context"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesl33/zk/internal/note"
)

// Yielder - TODO
type Yielder func(n note.Note, err error) bool

// List - TODO
func List(ctx context.Context, path string) iter.Seq2[note.Note, error] {
	return func(yield func(note.Note, error) bool) {
		err := filepath.WalkDir(path, func(path string, _ os.DirEntry, err error) error {
			return list(yield, path, err)
		})
		if err == nil {
			return
		}

		yield(note.Note{}, err)
	}
}

// list - TODO
func list(yield Yielder, path string, err error) error {
	// TODO
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if hidden(path) || !strings.HasSuffix(path, ".md") {
		return nil
	}

	// TODO
	if !yield(note.NewNote(path)) {
		return io.EOF
	}

	return nil
}
