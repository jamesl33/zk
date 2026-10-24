package vault

import (
	"fmt"
	"os"
	"path/filepath"
)

// Directories is the PARA layout created in every new vault.
var Directories = []string{
	"0 Inbox",
	"1 Projects",
	"2 Areas",
	"3 Resources",
	"4 Archives",
	"5 Bibliography",
}

// Init creates the vault (the '.zk' directory and PARA layout) in the given directory, which is created if required.
// Anything that exists is left untouched.
func Init(path string) error {
	for _, dir := range append([]string{".zk"}, Directories...) {
		err := os.MkdirAll(filepath.Join(path, dir), 0o755)
		if err != nil {
			return fmt.Errorf("failed to create '%s': %w", dir, err)
		}
	}

	return nil
}
