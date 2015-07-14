package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"straitjacket/engine"
	"strconv"
)

type executionStep struct {
	Stdout     string  `json:"stdout"`
	Stderr     string  `json:"stderr"`
	ExitStatus int     `json:"exit_status"`
	Time       float64 `json:"time"`
	Error      *string `json:"error"`
}

type executionResult struct {
	Success     bool           `json:"success"`
	Error       *string        `json:"error"`
	Compilation *executionStep `json:"compilation"`
	Runtime     *executionStep `json:"runtime"`
}

func (ctx *Context) ExecuteHandler(res http.ResponseWriter, req *http.Request) {
	languageName := req.FormValue("language")
	source := req.FormValue("source")
	stdin := req.FormValue("stdin")
	timelimit := req.FormValue("timelimit")

	log.Println(languageName, source, stdin, timelimit)

	language, err := ctx.Engine.FindLanguage(languageName)
	if err != nil {
		panic(err)
	}

	var timeout int64 = 60
	if timelimit != "" {
		timeout, err = strconv.ParseInt(timelimit, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	runResult, err := language.Run(&engine.RunOptions{
		Source:  source,
		Stdin:   stdin,
		Timeout: timeout,
	})
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(buildResult(runResult))
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}

func buildResult(runResult *engine.RunResult) *executionResult {
	result := &executionResult{}

	if runResult.RunStep != nil && runResult.RunStep.ErrorString != "" {
		result.Error = &runResult.RunStep.ErrorString
		result.Success = false
	} else if runResult.CompileStep != nil && runResult.CompileStep.ErrorString != "" {
		result.Error = &runResult.CompileStep.ErrorString
		result.Success = false
	} else {
		result.Success = true
	}

	if runResult.CompileStep != nil {
		result.Compilation = translateExecutionResult(runResult.CompileStep)
	}
	if runResult.RunStep != nil {
		result.Runtime = translateExecutionResult(runResult.RunStep)
	}

	return result
}

func translateExecutionResult(result *engine.ExecutionResult) *executionStep {
	res := &executionStep{
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
		ExitStatus: result.ExitCode,
		Time:       result.RunTime.Seconds(),
		Error:      &result.ErrorString,
	}

	if *res.Error == "" {
		res.Error = nil
	}

	return res
}
