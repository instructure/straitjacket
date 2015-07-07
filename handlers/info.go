package handlers

import (
	"encoding/json"
	"net/http"
)

type Option struct {
	VisibleName string `json:"visible_name"`
	Version     string `json:"version"`
}

type RuntimeOptions struct {
	Options map[string]Option `json:"languages"`
}

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	lang_map := make(map[string]Option)
	for _, lang := range ctx.Engine.Languages {
		lang_map[lang.Name] = Option{VisibleName: lang.VisibleName, Version: lang.Version}
	}

	options := RuntimeOptions{Options: lang_map}

	json, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}

	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}
