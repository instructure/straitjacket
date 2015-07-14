package handlers

import (
	"encoding/json"
	"net/http"
)

type Option struct {
	VisibleName string `json:"visible_name"`
	Version     string `json:"version"`
}

type AppInfo struct {
	Options    map[string]Option `json:"languages"`
	Extensions map[string]string `json:"extensions"`
}

var extensionsMap map[string]string

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	langMap := make(map[string]Option)
	for _, lang := range ctx.Engine.Languages {
		langMap[lang.Name] = Option{VisibleName: lang.VisibleName, Version: lang.Version}
	}

	if extensionsMap == nil {
		extensionsMap = make(map[string]string)
		for _, lang := range ctx.Engine.Languages {
			for _, ext := range lang.FileExtensions {
				extensionsMap[ext] = lang.Name
			}
		}
	}

	options := AppInfo{
		Options:    langMap,
		Extensions: extensionsMap,
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
