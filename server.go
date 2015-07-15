package main

import (
	"net/http"
	"straitjacket/engine"
	"straitjacket/handlers"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func newServerStack(engine *engine.Engine) *negroni.Negroni {
	context := &handlers.Context{
		Engine: engine,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", context.IndexHandler).Methods("GET")
	router.HandleFunc("/execute", context.ExecuteHandler).Methods("POST")
	router.HandleFunc("/info", context.InfoHandler).Methods("GET")

	c := cors.Default()
	server := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), c, negroni.NewStatic(http.Dir("public")))
	server.UseHandler(router)
	return server
}

func startServer(engine *engine.Engine, addr string) {
	server := newServerStack(engine)
	server.Run(addr)
}
