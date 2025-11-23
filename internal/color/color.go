package color

import "github.com/fatih/color"

func init() {
	color.NoColor = false
}

// Blue - TODO
var Blue = color.New(color.FgBlue).SprintFunc()

// Yellow - TODO
var Yellow = color.New(color.FgYellow).SprintFunc()
