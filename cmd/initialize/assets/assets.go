// Package assets contains the agent-agnostic files shared by every 'zk initialize' subcommand.
package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed instructions.md
var Instructions []byte

//go:embed skills
var Skills embed.FS

// CopySkills writes the embedded skills into 'dest' (e.g. '.claude/skills' or '.gemini/skills').
func CopySkills(dest string) error {
	return fs.WalkDir(Skills, "skills", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(path, "skills/")

		target := filepath.Join(dest, rel)

		err = os.MkdirAll(filepath.Dir(target), 0o755)
		if err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", target, err)
		}

		data, err := Skills.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		err = os.WriteFile(target, data, 0o644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", target, err)
		}

		return nil
	})
}
