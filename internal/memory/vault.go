package memory

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jamesl33/zk/cmd/initialize/assets"
	"github.com/jamesl33/zk/internal/vault"
)

// Guidance tells Claude how and when to save memories.
//
//go:embed assets/guidance.md
var Guidance string

// VaultPath returns the memory vault, from ZK_MEMORY_VAULT or '~/.claude/zk-memory'.
func VaultPath() (string, error) {
	if path := os.Getenv("ZK_MEMORY_VAULT"); path != "" {
		return filepath.Abs(path)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to find home directory: %w", err)
	}

	return filepath.Join(home, ".claude", "zk-memory"), nil
}

// EnsureVault creates the vault if required. It's idempotent, and the notes are the only thing git tracks: the index and
// hook state in '.zk' are disposable.
func EnsureVault(path string) error {
	err := vault.Init(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(stateDir(path), 0o755)
	if err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// The vault's own instructions, which the headless save run picks up as it runs in the vault.
	files := map[string]string{
		"ZK.md":     string(assets.Instructions),
		"CLAUDE.md": "@ZK.md\n\n" + Guidance,
	}

	for name, content := range files {
		err = writeIfChanged(filepath.Join(path, name), content)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(filepath.Join(path, ".git"))
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = writeIfChanged(filepath.Join(path, ".gitignore"), ".zk/\n")
	if err != nil {
		return err
	}

	out, err := exec.Command("git", "-C", path, "init", "-q").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize git: %w: %s", err, out)
	}

	return nil
}

// writeIfChanged writes the file unless it already holds the content.
func writeIfChanged(path, content string) error {
	existing, err := os.ReadFile(path)
	if err == nil && string(existing) == content {
		return nil
	}

	err = os.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '%s': %w", path, err)
	}

	return nil
}
