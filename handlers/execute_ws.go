package handlers

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func (ctx *Context) ExecuteWSHandler(res http.ResponseWriter, req *http.Request) {
	// ws, err := upgrader.Upgrade(res, req, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// defer ws.Close()
	//
	// var stdin bytes.Buffer
	// stdoutR, stdoutW := io.Pipe()
	// stderrR, stderrW := io.Pipe()
	// go readWS(ws, stdin)
	// writeWS(ws, stdoutR, stderrR)
	// _, err = ctx.Engine.Run(languageName, &engine.RunOptions{
	// 	Source:         source,
	// 	Stdin:          &stdin,
	// 	Stdout:         stdoutW,
	// 	Stderr:         stderrW,
	// 	Timeout:        timeout,
	// 	CompileTimeout: compileTimeout,
	// 	MaxOutputSize:  ctx.MaxOutputSize,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// stdoutW.Close()
	// stderrW.Close()
	return
}

func readWS(ws *websocket.Conn, stdin bytes.Buffer) {
	for {
		_, msg, err := ws.NextReader()
		if err != nil {
			break
		}
		io.Copy(&stdin, msg)
	}
}

type wsMessage struct {
	Stream, Message string
}

func writeWS(ws *websocket.Conn, stdout, stderr io.Reader) {
	// only one goroutine can write to the websocket at a time,
	// so we need synchronization of the two output streams.
	lock := &sync.Mutex{}
	go func() { io.Copy(&wsWriter{ws, lock, "stdout"}, stdout) }()
	go func() { io.Copy(&wsWriter{ws, lock, "stderr"}, stderr) }()
}

type wsWriter struct {
	ws *websocket.Conn
	*sync.Mutex
	stream string
}

func (wsWriter *wsWriter) Write(p []byte) (int, error) {
	message := struct {
		Stream, Message string
	}{wsWriter.stream, string(p)}
	wsWriter.Lock()
	defer wsWriter.Unlock()
	err := wsWriter.ws.WriteJSON(message)
	if err == nil {
		return len(p), nil
	}
	return 0, err
}
