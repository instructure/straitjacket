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

func (ctx *Context) InfoHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	langMap := make(map[string]Option)
	for _, lang := range ctx.Engine.Languages {
		langMap[lang.Name] = Option{VisibleName: lang.VisibleName, Version: lang.Version}
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

var extensionsMap = map[string]string{
	"sh":         "bash",
	"bash":       "bash",
	"c":          "c",
	"cs":         "c#",
	"c#":         "c#",
	"cpp":        "cpp",
	"cxx":        "cpp",
	"c++":        "cpp",
	"d":          "d",
	"f90":        "fortran",
	"fortran":    "fortran",
	"go":         "go",
	"guile":      "guile",
	"hs":         "haskell",
	"haskell":    "haskell",
	"java":       "java",
	"js":         "javascript",
	"sjs":        "javascript",
	"ssjs":       "javascript",
	"javascript": "javascript",
	"lua":        "lua",
	"ml":         "ocaml",
	"ocaml":      "ocaml",
	"pl":         "perl",
	"plx":        "perl",
	"perl":       "perl",
	"php":        "php",
	"php5":       "php",
	"py":         "python",
	"pyw":        "python",
	"xpy":        "python",
	"python":     "python",
	"rb":         "ruby",
	"rbw":        "ruby",
	"rbx":        "ruby",
	"ruby":       "ruby",
	"scala":      "scala",
	"rkt":        "scheme",
	"scm":        "scheme",
	"scheme":     "scheme",
}
