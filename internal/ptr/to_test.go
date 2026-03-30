package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	var (
		b   bool = true
		bp       = &b
		bpp      = To(bp)
	)

	// Asserts the value is the same
	assert.Equal(t, b, **bpp)

	// Asserts it's not a copy (other than a shallow copy)
	assert.Same(t, bp, *bpp)
}
