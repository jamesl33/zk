package color

import "github.com/fatih/color"

func init() {
	color.NoColor = false
}

// Yellow - TODO
var Yellow = color.New(color.FgYellow).SprintFunc()

// Blue - TODO
var Blue = color.New(color.FgBlue).SprintFunc()

// Cyan - TODO
var Cyan = color.New(color.FgCyan).SprintFunc()
