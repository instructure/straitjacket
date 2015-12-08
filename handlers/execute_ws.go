package handlers

import (
	"fmt"
	"io"
	"net/http"
	"straitjacket/engine"
	"strconv"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		sendResponse(status, &wsErrorResponse{StatusCode: strconv.Itoa(status), Error: reason.Error()}, w)
	},
}

type wsExecution struct {
	ID                       string
	output                   chan interface{}
	stdinR, stdoutR, stderrR io.Reader
	stdinW, stdoutW, stderrW io.WriteCloser
}

type executions struct {
	image      *engine.Image
	executions map[string]*wsExecution
	messages   chan interface{}
	wg         sync.WaitGroup
}

func (exes *executions) cleanup() {
	close(exes.messages)
	if exes.image != nil {
		exes.image.Remove()
	}
	exes.wg.Wait()
}

// ExecuteWSHandler handles websocket API connections.
func (ctx *Context) ExecuteWSHandler(res http.ResponseWriter, hreq *http.Request) {
	ws, err := upgrader.Upgrade(res, hreq, nil)
	if err != nil {
		// error response is written by the Upgrader callback defined above
		return
	}
	defer ws.Close()
	// 1024 fudge value is for the other input params like language
	ws.SetReadLimit((int64)(ctx.MaxSourceSize + ctx.MaxStdinSize + 1024))
	// SetReadTimeout ?

	exes := &executions{
		executions: map[string]*wsExecution{},
		messages:   make(chan interface{}),
	}
	exes.write(ctx.logger(hreq), ws)
	defer exes.cleanup()

	for {
		req := wsRequest{}
		err = ws.ReadJSON(&req)
		if _, ok := err.(*websocket.CloseError); ok {
			break
		}
		if err != nil {
			exes.messages <- wsError(nil, err)
			break
		}

		ctx.logger(hreq).WithFields(logrus.Fields{
			"request": req,
		}).Info("processing wsrequest")

		exes.process(ctx, &req)
	}
}

type wsProtocolError struct {
	msg string
}

func (e *wsProtocolError) Error() string {
	return e.msg
}

func wsError(id *string, err error) *wsErrorResponse {
	statusCode := "500"
	if _, ok := err.(*wsProtocolError); ok {
		statusCode = "400"
	}
	return &wsErrorResponse{
		ID:         id,
		StatusCode: statusCode,
		Error:      err.Error(),
	}
}

func (exes *executions) process(ctx *Context, req *wsRequest) {
	if req.Compile != nil {
		compileResult, err := exes.compile(ctx, req.Compile)
		if err == nil {
			exes.messages <- buildResult(&engine.RunResult{CompileStep: compileResult}, nil, nil)
		} else {
			exes.messages <- wsError(nil, err)
		}
	}
	if req.Run != nil {
		if exes.image == nil {
			exes.messages <- wsError(&req.Run.ID, &wsProtocolError{"must compile before running"})
		} else {
			execution := exes.newExecution(ctx, req.Run)
			execution.run(ctx, req.Run, exes.image)
		}
	}
	if req.Write != nil {
		execution := exes.executions[req.Write.ID]
		if execution == nil {
			exes.messages <- wsError(&req.Write.ID, &wsProtocolError{"specified id does not exist"})
		} else {
			execution.stdinW.Write([]byte(req.Write.Stdin))
		}
	}
	if req.Close != nil {
		execution := exes.executions[req.Close.ID]
		if execution == nil {
			exes.messages <- wsError(&req.Write.ID, &wsProtocolError{"specified id does not exist"})
		} else {
			execution.stdinW.Close()
		}
	}
}

func (exes *executions) write(logger *logrus.Entry, ws *websocket.Conn) {
	exes.wg.Add(1)
	go func() {
		defer exes.wg.Done()
		for msg := range exes.messages {
			logger.WithFields(logrus.Fields{"response": msg}).Info("response")
			if ws.WriteJSON(msg) != nil {
				break
			}
		}
	}()
}

func (exes *executions) compile(ctx *Context, req *wsCompileRequest) (*engine.ExecutionResult, error) {
	// remove any previous image
	if exes.image != nil {
		exes.image.Remove()
	}

	if req == nil || req.Language == "" || req.Source == "" {
		return nil, fmt.Errorf("required parameters are: language, source")
	}

	if req.Timeout == 0 {
		req.Timeout = ctx.DefaultTimeout
	}

	if len(req.Source) > ctx.MaxSourceSize {
		return nil, fmt.Errorf("source code is too large, max size: %d bytes", ctx.MaxSourceSize)
	}
	language := ctx.Engine.FindLanguage(req.Language)
	if language == nil {
		return nil, fmt.Errorf("unsupported language '%s'", req.Language)
	}

	image, res, err := language.Compile(req.Timeout, req.Source)
	exes.image = image
	return res, err
}

func (exes *executions) newExecution(ctx *Context, req *wsRunRequest) *wsExecution {
	if req.Timeout == 0 {
		req.Timeout = ctx.DefaultTimeout
	}
	exe := &wsExecution{ID: req.ID, output: exes.messages}
	exe.stdinR, exe.stdinW = io.Pipe()
	exe.stdoutR, exe.stdoutW = io.Pipe()
	exe.stderrR, exe.stderrW = io.Pipe()
	exes.executions[req.ID] = exe
	return exe
}

func (exe *wsExecution) run(ctx *Context, req *wsRunRequest, image *engine.Image) {
	go func() { io.Copy(&wsWriter{exe.ID, exe.output, "stdout"}, exe.stdoutR) }()
	go func() { io.Copy(&wsWriter{exe.ID, exe.output, "stderr"}, exe.stderrR) }()
	go func() {
		defer func() {
			exe.stdoutW.Close()
			exe.stderrW.Close()
		}()
		runStep, err := image.Run(&engine.RunOptions{
			Stdin:         exe.stdinR,
			Stdout:        exe.stdoutW,
			Stderr:        exe.stderrW,
			Timeout:       req.Timeout,
			MaxOutputSize: ctx.MaxOutputSize,
		})
		if err == nil {
			result := buildResult(&engine.RunResult{RunStep: runStep}, nil, nil)
			result.ID = exe.ID
			exe.output <- &result
		} else {
			exe.output <- wsError(&exe.ID, err)
		}
	}()
}

type wsWriter struct {
	ID     string
	output chan interface{}
	stream string
}

func (wsWriter *wsWriter) Write(p []byte) (int, error) {
	message := struct {
		ID     string `json:"id"`
		Stream string `json:"stream"`
		Output string `json:"output"`
	}{wsWriter.ID, wsWriter.stream, string(p)}
	wsWriter.output <- &message
	return len(p), nil
}
