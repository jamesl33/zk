package hs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfallible(t *testing.T) {
	var called bool

	fn := func(s string) {
		called = true
		assert.Equal(t, "test_input", s)
	}

	infallible := Infallible(fn)

	err := infallible("test_input")
	require.NoError(t, err)
	assert.True(t, called)
}
