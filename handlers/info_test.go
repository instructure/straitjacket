package handlers

import (
	"straitjacket/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

var source = []*engine.Language{
	{Name: "ruby", VisibleName: "Ruby 2.2", Version: "2.2", FileExtensions: []string{"rb"}},
	{Name: "d", VisibleName: "D (GDC)", Version: "2.5.5.9", FileExtensions: []string{"d", "dd"}},
}

func TestLangMap(t *testing.T) {
	langs := langMap(source)
	assert.Equal(t, 2, len(langs))
	assert.Equal(t, "Ruby 2.2", langs["ruby"].VisibleName)
	assert.Equal(t, "D (GDC)", langs["d"].VisibleName)
	assert.Equal(t, "2.2", langs["ruby"].Version)
	assert.Equal(t, "2.5.5.9", langs["d"].Version)
}

func TestExtensionMap(t *testing.T) {
	exts := makeExtensionsMap(source)
	assert.Equal(t, 3, len(exts))
	assert.Equal(t, "ruby", exts["rb"])
	assert.Equal(t, "d", exts["d"])
	assert.Equal(t, "d", exts["dd"])
}
