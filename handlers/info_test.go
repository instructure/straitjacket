package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLangMap(t *testing.T) {
	langs := langMap(testLangs)
	assert.Equal(t, 2, len(langs))
	assert.Equal(t, "Ruby 2.2", langs["ruby"].VisibleName)
	assert.Equal(t, "D (GDC)", langs["d"].VisibleName)
	assert.Equal(t, "2.2", langs["ruby"].Version)
	assert.Equal(t, "2.5.5.9", langs["d"].Version)
}

func TestExtensionMap(t *testing.T) {
	exts := makeExtensionsMap(testLangs)
	assert.Equal(t, 3, len(exts))
	assert.Equal(t, "ruby", exts["rb"])
	assert.Equal(t, "d", exts["d"])
	assert.Equal(t, "d", exts["dd"])
}

func TestInfoResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	engine := NewMockEngine(mockCtrl)
	engine.EXPECT().Languages().AnyTimes().Return(testLangs)

	req, _ := http.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()
	ctx := &Context{
		Engine: engine,
	}
	ctx.InfoHandler(w, req)

	expected := `{
	  "languages": {
		  "ruby": { "visible_name": "Ruby 2.2", "version": "2.2" },
			"d": { "visible_name": "D (GDC)", "version": "2.5.5.9" }
		},
		"extensions": {
		  "rb": "ruby",
			"d": "d",
			"dd": "d"
		}
	}`

	assertJSONResponse(t, expected, w)
}
