package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/jamesl33/zk/internal/vector"
	"github.com/spf13/cobra"
)

// FindOptions defines the options for the find command.
type FindOptions struct {
	// Limit is the maximum number of notes to find.
	Limit int

	// MaxDistance is the maximum cosine distance a note may be from the query.
	MaxDistance float32

	// JSON outputs each note as a JSON line.
	JSON bool
}

// Find defines the struct for the find command.
type Find struct {
	FindOptions
}

// NewFind creates a new command for finding semantically similar notes.
func NewFind() *cobra.Command {
	var find Find

	cmd := cobra.Command{
		Short: "Finds semantically similar notes using a plain text query",
		Use:   "find <query>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().IntVar(&find.Limit, "limit", 0, "Maximum number of notes to find (default 15)")
	cmd.Flags().Float32Var(&find.MaxDistance, "max-distance", 0, "Maximum cosine distance from the query (default unbounded)")
	cmd.Flags().BoolVar(&find.JSON, "json", false, "Output each note as a JSON line")

	return &cmd
}

// Run finds some related notes.
func (f *Find) Run(ctx context.Context, query string) error {
	if _, err := vault.RootAbs("."); err != nil {
		return err
	}

	n := &note.Note{
		// We just want to set the body
	}

	n.SetBody(query)

	opts := []func(*vector.FindOptions){
		vector.WithLimit(f.Limit),
		vector.WithMaxDistance(f.MaxDistance),
	}

	return notes.Find(ctx, n, f.output, opts...)
}

// output prints the note, as a JSON line if requested.
func (f *Find) output(n *note.Note) error {
	if f.JSON {
		return json.NewEncoder(os.Stdout).Encode(n)
	}

	fmt.Println(n.String())

	return nil
}
