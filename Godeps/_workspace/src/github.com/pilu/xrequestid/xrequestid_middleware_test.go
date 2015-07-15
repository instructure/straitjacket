package xrequestid

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXRequestID(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	middleware := New(16)
	middleware.Generate = func(n int) (string, error) { return "test-id", nil }
	middleware.ServeHTTP(recorder, req, func(w http.ResponseWriter, r *http.Request) {})

	if id := req.Header.Get("X-Request-ID"); id != "test-id" {
		t.Fatalf("Expected X-Request-Id to be `test-id`, got `%v`", id)
	}
}
