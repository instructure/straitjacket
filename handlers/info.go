package handlers

import (
	"encoding/json"
	"net/http"
)

type Option struct {
	VisibleName      string `json:"visible_name"`
	Version          string `json:"version"`
	ExecutionProfile string `json:"execution_profile"`
}

type RuntimeOptions struct {
	Options []Option `json:"languages"`
}

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// TODO: read config directory for language options
	options := RuntimeOptions{
		[]Option{
			Option{VisibleName: "JavaScript", Version: "0.12", ExecutionProfile: "javascript"},
			Option{VisibleName: "Ruby", Version: "2.2", ExecutionProfile: "ruby"},
			Option{VisibleName: "C#", Version: "4.0", ExecutionProfile: "csharp"},
		},
	}

	json, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}

	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}
