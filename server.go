package main

import (
	"log"
	"net/http"
	"straitjacket/engine"
	"straitjacket/handlers"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/pilu/xrequestid"
	"github.com/rs/cors"
)

func newServerStack(engine *engine.Engine) *negroni.Negroni {
	context := handlers.NewContext(engine)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", context.IndexHandler).Methods("GET")
	router.HandleFunc("/execute", context.ExecuteHandler).Methods("POST")
	router.HandleFunc("/executews", context.ExecuteWSHandler).Methods("GET")
	router.HandleFunc("/info", context.InfoHandler).Methods("GET")

	server := negroni.New(xrequestid.New(16),
		context,
		negronilogrus.NewCustomMiddleware(logrus.InfoLevel, &logrus.JSONFormatter{}, "straitjacket"),
		cors.Default(),
		negroni.NewStatic(http.Dir("public")))
	server.UseHandler(router)
	return server
}

func startServer(engine *engine.Engine, addr string) {
	server := newServerStack(engine)
	log.Fatal(http.ListenAndServe(addr, server))
}
