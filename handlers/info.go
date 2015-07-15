package handlers

import (
	"encoding/json"
	"net/http"
	"straitjacket/engine"
)

type language struct {
	VisibleName string `json:"visible_name"`
	Version     string `json:"version"`
}

type appInfo struct {
	Languages  map[string]language `json:"languages"`
	Extensions map[string]string   `json:"extensions"`
}

var extensionsMap map[string]string

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if extensionsMap == nil {
		extensionsMap = makeExtensionsMap(ctx.Engine.Languages)
	}

	options := appInfo{
		Languages:  langMap(ctx.Engine.Languages),
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

func langMap(languages []*engine.Language) map[string]language {
	langMap := make(map[string]language)
	for _, lang := range languages {
		langMap[lang.Name] = language{VisibleName: lang.VisibleName, Version: lang.Version}
	}
	return langMap
}

func makeExtensionsMap(languages []*engine.Language) map[string]string {
	extensionsMap = make(map[string]string)
	for _, lang := range languages {
		for _, ext := range lang.FileExtensions {
			extensionsMap[ext] = lang.Name
		}
	}
	return extensionsMap
}
