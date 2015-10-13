// Package xrequestid implements an http middleware for Negroni that assigns a random id to each request
//
// Example:
//		package main
//
//		import (
//			"fmt"
//			"net/http"
//
//			"github.com/codegangsta/negroni"
//			"github.com/pilu/xrequestid"
//		)
//
//		func main() {
//			mux := http.NewServeMux()
//			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//				fmt.Fprintf(w, "X-Request-Id is `%s`", r.Header.Get("X-Request-Id"))
//			})
//
//			n := negroni.New()
//			n.Use(xrequestid.New(16))
//			n.UseHandler(mux)
//			n.Run(":3000")
//		}
package xrequestid

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// By default the middleware set the generated random string to this key in the request header
const DefaultHeaderKey = "X-Request-Id"

// GenerateFunc is the func used by the middleware to generates the random string.
type GenerateFunc func(int) (string, error)

// XRequestID is a middleware that adds a random ID to the request X-Request-Id header
type XRequestID struct {
	// Size specifies the length of the random length. The length of the result string is twice of n.
	Size int
	// Generate is a GenerateFunc that generates the random string. The default one uses crypto/rand
	Generate GenerateFunc
	// HeaderKey is the header name where the middleware set the random string. By default it uses the DefaultHeaderKey constant value
	HeaderKey string
}

// New returns a new XRequestID middleware instance. n specifies the length of the random length. The length of the result string is twice of n.
func New(n int) *XRequestID {
	return &XRequestID{
		Size:      n,
		Generate:  generateID,
		HeaderKey: DefaultHeaderKey,
	}
}

func (m *XRequestID) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := m.Generate(m.Size)
	if err == nil {
		r.Header.Set(m.HeaderKey, id)
	}

	next(rw, r)
}

func generateID(n int) (string, error) {
	r := make([]byte, n)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(r), nil
}
