package handlers

import "straitjacket/engine"

type Engine interface {
	Languages() []*engine.Language
	Run(languageName string, opts *engine.RunOptions) (*engine.RunResult, error)
}

type Context struct {
	Engine        Engine
	extensionsMap map[string]string
}
