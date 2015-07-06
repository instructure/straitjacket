package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"straitjacket/handlers"
)

func NewServerStack() *negroni.Negroni {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/execute", handlers.ExecuteHandler)
	router.HandleFunc("/info", handlers.InfoHandler)

	server := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	server.UseHandler(router)
	return server
}

func StartServer(addr string) {
	server := NewServerStack()
	server.Run(addr)
}
