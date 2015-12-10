package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"straitjacket/engine"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type executionStep struct {
	Stdout     string  `json:"stdout,omitempty"`
	Stderr     string  `json:"stderr,omitempty"`
	ExitStatus int     `json:"exit_status"`
	Time       float64 `json:"time"`
	Error      *string `json:"error"`
}

type executionResult struct {
	Success     bool           `json:"success"`
	Error       *string        `json:"error"`
	Compilation *executionStep `json:"compilation,omitempty"`
	Runtime     *executionStep `json:"runtime,omitempty"`
	// used by the websocket API responses
	ID         string `json:"id,omitempty"`
	StatusCode string `json:"status_code,omitempty"`
}

func (ctx *Context) ExecuteHandler(res http.ResponseWriter, req *http.Request) {
	languageName := req.FormValue("language")
	source := req.FormValue("source")
	stdin := req.FormValue("stdin")
	timelimit := req.FormValue("timelimit")
	compileTimelimit := req.FormValue("compile_timelimit")

	if len(source) > ctx.MaxSourceSize || len(stdin) > ctx.MaxStdinSize {
		errorResponse(413, "request_size_error", res)
		return
	}

	ctx.logger(req).WithFields(logrus.Fields{
		"language":          languageName,
		"source":            source,
		"stdin":             stdin,
		"timelimit":         timelimit,
		"compile_timelimit": compileTimelimit,
	}).Info("executing code")

	timeout := parseTimelimit(timelimit, ctx.DefaultTimeout)
	compileTimeout := parseTimelimit(compileTimelimit, timeout)

	var stdout, stderr bytes.Buffer
	lang := ctx.Engine.FindLanguage(languageName)
	if lang == nil {
		panic(fmt.Errorf("Language not found: '%s'", languageName))
	}

	runResult, err := lang.Run(&engine.RunOptions{
		Source:         source,
		Stdin:          strings.NewReader(stdin),
		Stdout:         &stdout,
		Stderr:         &stderr,
		Timeout:        timeout,
		CompileTimeout: compileTimeout,
		MaxOutputSize:  ctx.MaxOutputSize,
	})
	if err != nil {
		panic(err)
	}

	response := buildResult(runResult, &stdout, &stderr)
	code := 200
	if !response.Success {
		code = 400
	}
	sendResponse(code, response, res)
}

func errorResponse(code int, message string, res http.ResponseWriter) {
	response := &executionResult{
		Success: false,
		Error:   &message,
	}
	sendResponse(code, response, res)
}

func sendResponse(code int, response interface{}, res http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	_, err = res.Write(json)
	if err != nil {
		panic(err)
	}
}

func parseTimelimit(timelimit string, defaultLimit int64) (timeout int64) {
	var err error
	if timelimit != "" {
		timeout, err = strconv.ParseInt(timelimit, 10, 64)
	}
	if err != nil || timelimit == "" {
		timeout = defaultLimit
	}
	return
}

func buildResult(runResult *engine.RunResult, stdout, stderr *bytes.Buffer) *executionResult {
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
		result.Compilation = translateExecutionResult(runResult.CompileStep, nil, nil)
	}
	if runResult.RunStep != nil {
		result.Runtime = translateExecutionResult(runResult.RunStep, stdout, stderr)
	}

	return result
}

func translateExecutionResult(result *engine.ExecutionResult, stdout, stderr *bytes.Buffer) *executionStep {
	var stdoutStr, stderrStr string

	if stdout == nil {
		stdoutStr = result.Stdout
	} else {
		stdoutStr = stdout.String()
	}

	if stderr == nil {
		stderrStr = result.Stderr
	} else {
		stderrStr = stderr.String()
	}

	res := &executionStep{
		Stdout:     stdoutStr,
		Stderr:     stderrStr,
		ExitStatus: result.ExitCode,
		Time:       result.RunTime.Seconds(),
		Error:      &result.ErrorString,
	}

	if *res.Error == "" {
		res.Error = nil
	}

	return res
}
