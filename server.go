package main

import (
	"straitjacket/engine"
	"straitjacket/handlers"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func newServerStack(engine *engine.Engine) *negroni.Negroni {
	context := &handlers.Context{
		Engine: *engine,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", context.IndexHandler)
	router.HandleFunc("/execute", context.ExecuteHandler)
	router.HandleFunc("/info", context.InfoHandler)

	server := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	server.UseHandler(router)
	return server
}

func startServer(engine *engine.Engine, addr string) {
	server := newServerStack(engine)
	server.Run(addr)
}
