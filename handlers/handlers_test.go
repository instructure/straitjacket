package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"straitjacket/engine"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertJSONResponse(t *testing.T, expectedResponse string, w *httptest.ResponseRecorder) {
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("content-type"))

	var expected, actual interface{}
	require.NoError(t, json.Unmarshal([]byte(expectedResponse), &expected))
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &actual))

	assert.Equal(t, expected, actual)
}

var testLangs = []*engine.Language{
	{Name: "ruby", VisibleName: "Ruby 2.2", Version: "2.2", FileExtensions: []string{"rb"}},
	{Name: "d", VisibleName: "D (GDC)", Version: "2.5.5.9", FileExtensions: []string{"d", "dd"}},
}
