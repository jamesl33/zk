// Package memory implements the Claude Code hooks which back Claude's persistent memory with a zk vault.
package memory

import (
	"encoding/json"
	"io"
	"os"
)

// SavingEnv is set for the headless run which saves memories, so that its own hooks do nothing.
const SavingEnv = "ZK_MEMORY_SAVING"

// Saving returns whether this process is part of a headless save run.
func Saving() bool {
	return os.Getenv(SavingEnv) != ""
}

// Input is the JSON Claude Code writes to a hook's stdin.
type Input struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
	Source         string `json:"source"`
	Prompt         string `json:"prompt"`
}

// ReadInput decodes a hook's input.
func ReadInput(r io.Reader) (Input, error) {
	var in Input

	err := json.NewDecoder(r).Decode(&in)

	return in, err
}

// Output is the JSON a hook writes to stdout.
type Output struct {
	SystemMessage      string    `json:"systemMessage,omitempty"`
	HookSpecificOutput *Specific `json:"hookSpecificOutput,omitempty"`
}

// Specific is the event specific part of a hook's output.
type Specific struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext,omitempty"`
}

// Context returns output adding the given context for the event.
func Context(event, context string) Output {
	specific := Specific{HookEventName: event, AdditionalContext: context}

	out := Output{HookSpecificOutput: &specific}

	return out
}
