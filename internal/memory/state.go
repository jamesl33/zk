package memory

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// stateMaxAge is how long a session's state is kept after it was last touched.
	stateMaxAge = 7 * 24 * time.Hour

	// logMaxSize is the size at which the log is started afresh.
	logMaxSize = 256 << 10
)

// stateDir holds everything which is disposable: the index lives beside it in '.zk'.
func stateDir(vault string) string {
	return filepath.Join(vault, ".zk", "state")
}

// sessionFile names a file of per-session state, ensuring the session id can't escape the directory.
func sessionFile(vault, session, ext string) string {
	return filepath.Join(stateDir(vault), filepath.Base(session)+"."+ext)
}

// Recalled returns the notes already injected into the session.
func Recalled(vault, session string) (map[string]bool, error) {
	file, err := os.Open(sessionFile(vault, session, "recalled"))
	if errors.Is(err, os.ErrNotExist) {
		return map[string]bool{}, nil
	}

	if err != nil {
		return nil, err // Purposefully not wrapped
	}
	defer file.Close()

	recalled := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		recalled[scanner.Text()] = true
	}

	return recalled, scanner.Err() // Purposefully not wrapped
}

// MarkRecalled records notes as injected into the session.
func MarkRecalled(vault, session string, paths ...string) error {
	file, err := os.OpenFile(sessionFile(vault, session, "recalled"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err // Purposefully not wrapped
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(paths, "\n") + "\n")

	return err // Purposefully not wrapped
}

// ResetRecalled forgets what was injected into the session, as compaction drops it from context.
func ResetRecalled(vault, session string) error {
	err := os.Remove(sessionFile(vault, session, "recalled"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err // Purposefully not wrapped
}

// Offset returns the number of transcript lines already extracted for the session.
func Offset(vault, session string) int {
	data, err := os.ReadFile(sessionFile(vault, session, "offset"))
	if err != nil {
		return 0
	}

	offset, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}

	return offset
}

// SetOffset records the number of transcript lines extracted for the session.
func SetOffset(vault, session string, offset int) error {
	err := os.WriteFile(sessionFile(vault, session, "offset"), []byte(strconv.Itoa(offset)), 0o644)

	return err // Purposefully not wrapped
}

// Prune removes state which hasn't been touched recently. The log is kept.
func Prune(vault string) error {
	entries, err := os.ReadDir(stateDir(vault))
	if err != nil {
		return fmt.Errorf("failed to read state directory: %w", err)
	}

	for _, entry := range entries {
		switch filepath.Ext(entry.Name()) {
		case ".recalled", ".offset":
		default:
			continue
		}

		info, err := entry.Info()
		if err != nil || time.Since(info.ModTime()) < stateMaxAge {
			continue
		}

		_ = os.Remove(filepath.Join(stateDir(vault), entry.Name()))
	}

	return nil
}

// Logf appends to the hooks log. Hooks fail open, so this is the only trace of a failure.
func Logf(vault, format string, args ...any) {
	path := filepath.Join(stateDir(vault), "hooks.log")

	info, err := os.Stat(path)
	if err == nil && info.Size() > logMaxSize {
		_ = os.Remove(path)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "%s %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
}
