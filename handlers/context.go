package handlers

import (
	"net/http"
	"straitjacket/engine"

	"github.com/Sirupsen/logrus"
)

// Engine is an execution engine interface abstraction, primarily to aid in
// mocking during tests. Currently engine.Engine is the only concrete implementation.
type Engine interface {
	Languages() []*engine.Language
	Run(languageName string, opts *engine.RunOptions) (*engine.RunResult, error)
}

// Context is the HTTP handler context.
type Context struct {
	Engine        Engine
	extensionsMap map[string]string
	log           *logrus.Logger
}

// NewContext returns a new HTTP handler context that will use the provided
// execution engine.
func NewContext(engine Engine) *Context {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.JSONFormatter{}
	return &Context{
		Engine: engine,
		log:    log,
	}
}

func (ctx *Context) logger(req *http.Request) *logrus.Entry {
	log := logrus.NewEntry(ctx.log)
	if reqID := req.Header.Get("X-Request-Id"); reqID != "" {
		log = log.WithField("request_id", reqID)
	}
	return log
}
