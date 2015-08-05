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
	Template    string `json:"template"`
}

var infoResponse struct {
	Languages  []*language       `json:"languages"`
	Extensions map[string]string `json:"extensions"`
}

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	if infoResponse.Languages == nil {
		makeInfoResponse(ctx.Engine.Languages())
	}

	json, err := json.Marshal(infoResponse)
	if err != nil {
		panic(err)
	}

	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}

func makeInfoResponse(languages []*engine.Language) {
	infoResponse.Languages = langList(languages)
	infoResponse.Extensions = makeExtensionsMap(languages)
}

func langList(languages []*engine.Language) (langList []*language) {
	for _, lang := range languages {
		info := &language{
			Name:        lang.Name,
			VisibleName: lang.VisibleName,
			Version:     lang.Version,
			Template:    lang.Checks.Simple.Source,
		}
		langList = append(langList, info)
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
