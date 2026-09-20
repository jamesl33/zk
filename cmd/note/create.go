package note

import (
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/spf13/cobra"
)

// NewCreate creates a new command for creating notes.
func NewCreate() *cobra.Command {
	cmd := cobra.Command{
		Short: "Create a new note",
		Use:   "create",
	}

	cmd.AddCommand(
		NewCreateBibliographic(),
		NewCreatePermanent(),
		NewCreateFleeting(),
		NewCreateIndex(),
		NewCreateLiterature(),
	)

	return &cmd
}

// create writes a new note of the given type/title at the given path and
// prints the resulting path.
func create(noteType, title, path string) error {
	if _, err := vault.Root("."); err != nil {
		return err
	}

	fm := note.Frontmatter{
		Type:  note.Type(noteType),
		Title: title,
		Date:  time.Now().Format("2006-01-02"),
		Tags:  make([]string, 0),
	}

	n := note.Note{
		Path:        note.Path(path),
		Frontmatter: fm,
	}

	err := n.Create()
	if err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}

	fmt.Printf("%s\n", n.Path)

	return nil
}
