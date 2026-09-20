package regex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkMatch(t *testing.T) {
	m := Link.FindStringSubmatch("[[20060102150405]]")
	require := Link.SubexpNames()

	assert.NotNil(t, m)
	assert.Equal(t, "20060102150405", m[indexOf(require, "link")])
	assert.Empty(t, m[indexOf(require, "text")])
}

func TestLinkMatchWithText(t *testing.T) {
	m := Link.FindStringSubmatch("[[20060102150405|Some Text]]")
	names := Link.SubexpNames()

	assert.NotNil(t, m)
	assert.Equal(t, "20060102150405", m[indexOf(names, "link")])
	assert.Equal(t, "Some Text", m[indexOf(names, "text")])
}

func TestLinkNoMatch(t *testing.T) {
	m := Link.FindStringSubmatch("not a link")
	assert.Nil(t, m)
}

func indexOf(names []string, name string) int {
	for i, n := range names {
		if n == name {
			return i
		}
	}

	return -1
}
