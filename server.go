package main

import (
	"straitjacket/engine"
	"straitjacket/handlers"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func newServerStack() *negroni.Negroni {
	theEngine, err := engine.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	context := &handlers.Context{
		Engine: theEngine,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", context.IndexHandler)
	router.HandleFunc("/execute", context.ExecuteHandler)
	router.HandleFunc("/info", context.InfoHandler)

	server := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	server.UseHandler(router)
	return server
}

func startServer(addr string) {
	server := newServerStack()
	server.Run(addr)
}
