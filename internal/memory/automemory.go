package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AutoMemoryEnabled returns whether Claude Code's built-in auto-memory is on, which it is unless the environment or a
// settings file turns it off. A plugin can't turn it off itself, so the user has to.
func AutoMemoryEnabled(cwd string) bool {
	if os.Getenv("CLAUDE_CODE_DISABLE_AUTO_MEMORY") == "1" {
		return false
	}

	home, _ := os.UserHomeDir()

	for _, path := range []string{
		filepath.Join(home, ".claude", "settings.json"),
		filepath.Join(cwd, ".claude", "settings.json"),
		filepath.Join(cwd, ".claude", "settings.local.json"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var settings struct {
			AutoMemoryEnabled *bool `json:"autoMemoryEnabled"`
		}

		if json.Unmarshal(data, &settings) == nil && settings.AutoMemoryEnabled != nil && !*settings.AutoMemoryEnabled {
			return false
		}
	}

	return true
}
