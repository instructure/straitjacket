package handlers

import (
	"encoding/json"
	"net/http"
	"straitjacket/engine"
)

type language struct {
	Name        string `json:"name"`
	VisibleName string `json:"visible_name"`
	Version     string `json:"version"`
}

type appInfo struct {
	Languages  []*language       `json:"languages"`
	Extensions map[string]string `json:"extensions"`
}

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if ctx.extensionsMap == nil {
		ctx.extensionsMap = makeExtensionsMap(ctx.Engine.Languages())
	}

	options := appInfo{
		Languages:  langList(ctx.Engine.Languages()),
		Extensions: ctx.extensionsMap,
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

func langList(languages []*engine.Language) (langList []*language) {
	for _, lang := range languages {
		langList = append(langList, &language{Name: lang.Name, VisibleName: lang.VisibleName, Version: lang.Version})
	}
	return
}

func makeExtensionsMap(languages []*engine.Language) map[string]string {
	extensionsMap := make(map[string]string)
	for _, lang := range languages {
		for _, ext := range lang.FileExtensions {
			extensionsMap[ext] = lang.Name
		}
	}
	return extensionsMap
}
