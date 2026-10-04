// Package assets contains the files shared by every 'zk initialize' subcommand.
package assets

import (
	_ "embed"
)

//go:embed instructions.md
var Instructions []byte
