package notes

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

// FindOptions - TODO
type FindOptions struct {
	// Similar - TODO
	Similar bool

	// Dissimilar - TODO
	Dissimilar bool
}

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
		RunE: func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&find.Similar,
		"similar",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&find.Dissimilar,
		"dissimilar",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (f *Find) Run(ctx context.Context, args []string) error {
	// path - TODO
	const path = "zk.sqlite3"

	db, err := sql.Open("sqlite3", filepath.Join(os.TempDir(), path)+"?mode=ro")
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer db.Close()

	// TODO (jamesl33): Pass these arguments properly.
	n, err := note.New(args[0])
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	embedding, err := f.embed(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if embedding == nil {
		return nil
	}

	serial, err := sqlite_vec.SerializeFloat32(embedding)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// query - TODO
	const query = `
	SELECT
	  name,
	  checksum,
	  vec_distance_cosine (embedding, ?) AS distance
	FROM
	  notes
	WHERE
	  name != ? AND distance <= 0.6
	ORDER BY
	  distance
	LIMIT
	  25;
	`

	rows, err := db.QueryContext(ctx, query, serial, n.Name())
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	var (
		name     string
		checksum uint32
		distance float32
	)

	// TODO
	names := make([]string, 0)

	for rows.Next() {
		err := rows.Scan(&name, &checksum, &distance)
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(names) == 0 {
		return nil
	}

	err = f.list(ctx, names...)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// embed - TODO
//
// TODO (jamesl33): Create a vector search package.
func (f *Find) embed(ctx context.Context, n *note.Note) ([]float32, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	req := api.EmbedRequest{
		Model: "embeddinggemma:300m",
		Input: n.Body,
	}

	resp, err := client.Embed(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	return resp.Embeddings[0], nil
}

// list - TODO
func (f *Find) list(ctx context.Context, names ...string) error {
	matchers := hs.Map(names, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fmt.Println(n.String0())
	}

	return nil
}
