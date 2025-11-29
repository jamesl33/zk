package notes

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/database/vector"
	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// FindOptions - TODO
type FindOptions struct{}

// Find - TODO
type Find struct {
	FindOptions
}

// NewFind - TODO
func NewFind() *cobra.Command {
	var find Find

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "find",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

// Run - TODO
func (f *Find) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer db.Close()

	notes, err := db.Find(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(notes, hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
