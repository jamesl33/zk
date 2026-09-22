package ai

import "errors"

// ErrExceededContextLength is returned when the input exceeds the model's context length.
var ErrExceededContextLength = errors.New("input exceeds the context length")
