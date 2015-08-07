package engine

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteLimit(t *testing.T) {
	var raw bytes.Buffer
	buf := NewLimitedWriter(&raw, 5)

	sz, err := buf.Write(([]byte)("ab"))
	assert.Equal(t, 2, sz)
	assert.Nil(t, err)

	sz, err = buf.Write(([]byte)("ab"))
	assert.Equal(t, 2, sz)
	assert.Nil(t, err)

	sz, err = buf.Write(([]byte)("ab"))
	assert.Equal(t, 1, sz)
	assert.IsType(t, &OutputTooLarge{}, err)

	sz, err = buf.Write(([]byte)("ab"))
	assert.Equal(t, 0, sz)
	assert.IsType(t, &OutputTooLarge{}, err)

	assert.Equal(t, "ababa", raw.String())
}

func TestZeroWrite(t *testing.T) {
	var raw bytes.Buffer
	buf := NewLimitedWriter(&raw, 2)

	sz, err := buf.Write(([]byte)(""))
	assert.Equal(t, 0, sz)
	assert.Nil(t, err)

	buf.Write(([]byte)("ab"))

	// shouldn't error on zero-length writes even when full
	sz, err = buf.Write(([]byte)(""))
	assert.Equal(t, 0, sz)
	assert.Nil(t, err)
}
