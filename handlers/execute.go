package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ExecutionResult struct {
	STDOut     string `json:"stdout"`
	STDErr     string `json:"stderr"`
	ExitStatus int    `json:"exit_status"`
	Time       string `json:"time"`
	Error      string `json:"error"`
}

func ExecuteHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	lang := req.PostForm["language"]
	source := req.PostForm["source"]
	stdin := req.PostForm["stdin"]
	timelimit := req.PostForm["timelimit"]

	log.Println(lang, source, stdin, timelimit)

	result := ExecutionResult{
		STDOut:     "Fake Return",
		ExitStatus: 0,
	}
	json, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}
