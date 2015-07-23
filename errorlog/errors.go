package errorlog

import (
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
)

// Error logs an error event to the external error tracker, and returns a unique
// identifier for the event.
func Error(err error, tags map[string]string) string {
	return raven.CaptureError(err, tags)
}

// HTTPError logs an error event in the context of a http request, and returns a
// unique identifier for the event.
func HTTPError(err error, req *http.Request, tags map[string]string) string {
	return raven.CaptureError(err, tags, raven.NewHttp(req))
}

func init() {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		err := raven.SetDSN(dsn)
		if err != nil {
			panic(err)
		}
	}
}
