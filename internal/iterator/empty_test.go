package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	var (
		it    = Empty[string]()
		count = 0
	)

	for range it {
		count++
	}

	assert.Zero(t, count)
}
