package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"straitjacket/engine"
	"straitjacket/errorlog"

	"github.com/Sirupsen/logrus"
)

// Engine is an execution engine interface abstraction, primarily to aid in
// mocking during tests. Currently engine.Engine is the only concrete implementation.
type Engine interface {
	Languages() []*engine.Language
	FindLanguage(languageName string) *engine.Language
}

// Context is the HTTP handler context.
type Context struct {
	Engine         Engine
	extensionsMap  map[string]string
	log            *logrus.Logger
	DefaultTimeout int64
	MaxSourceSize  int
	MaxStdinSize   int
	MaxOutputSize  int
}

// NewContext returns a new HTTP handler context that will use the provided
// execution engine.
func NewContext(engine Engine) *Context {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.JSONFormatter{}
	return &Context{
		Engine:         engine,
		log:            log,
		DefaultTimeout: 60,
		MaxSourceSize:  64 * 1024,
		MaxStdinSize:   512 * 1024,
		MaxOutputSize:  64 * 1024,
	}
}

func (ctx *Context) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]

			errorid := errorlog.HTTPError(fmt.Errorf("%v", err), r, nil)

			ctx.logger(r).
				WithField("stack", string(stack)).
				WithField("error_id", errorid).
				Errorf("panic: %v\n", err)
		}
	}()

	// 1024 fudge value is for the other input params like language
	r.Body = http.MaxBytesReader(rw, r.Body, (int64)(ctx.MaxSourceSize+ctx.MaxStdinSize+1024))

	next(rw, r)
}

func requestID(req *http.Request) string {
	return req.Header.Get("X-Request-Id")
}

func (ctx *Context) logger(req *http.Request) *logrus.Entry {
	log := logrus.NewEntry(ctx.log)
	if reqID := requestID(req); reqID != "" {
		log = log.WithField("request_id", reqID)
	}
	return log
}
