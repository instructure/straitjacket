package handlers

import (
	"encoding/json"
	"net/http"
	"straitjacket/engine"
	"strconv"

	"github.com/Sirupsen/logrus"
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

	ctx.logger(req).WithFields(logrus.Fields{
		"language":  languageName,
		"source":    source,
		"stdin":     stdin,
		"timelimit": timelimit,
	}).Info("executing code")

	timeout, err := parseTimelimit(timelimit)
	if err != nil {
		panic(err)
	}

	runResult, err := ctx.Engine.Run(languageName, &engine.RunOptions{
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

func parseTimelimit(timelimit string) (timeout int64, err error) {
	timeout = 60
	if timelimit != "" {
		timeout, err = strconv.ParseInt(timelimit, 10, 64)
	}
	return
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
