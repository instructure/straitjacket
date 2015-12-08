package handlers

type wsRequest struct {
	Compile *wsCompileRequest `json:"compile,omitempty"`
	Run     *wsRunRequest     `json:"run,omitempty"`
	Write   *wsStdinRequest   `json:"write,omitempty"`
	Close   *wsCloseRequest   `json:"close,omitempty"`
}

// {"run": {"id": "1"}, "write": {"id": "1", "stdin": "1\n2\n3\n"}, "close": {"id": "1"}}

type wsCompileRequest struct {
	Language string `json:"language"`
	Source   string `json:"source"`
	Timeout  int64  `json:"timelimit,omitempty"`
}

// {"compile": {"language": "ruby", "source": "$stdout.sync = true; $stdin.each_line { |l| puts(l.to_i + 1) }"}}
// {"compile": {"language": "cpp", "source": "#include <iostream>\nint main() { std::cout << \"hey\\n\"; return 0; }"}}

type wsRunRequest struct {
	ID      string `json:"id"`
	Timeout int64  `json:"timelimit,omitempty"`
}

// {"run": {"id": "1"}}

type wsStdinRequest struct {
	ID    string `json:"id"`
	Stdin string `json:"stdin"`
}

// {"write": {"id": "1", "stdin": "1\n2\n3\n"}}

type wsCloseRequest struct {
	ID string `json:"id"`
}

// {"close": {"id": "1"}}

type wsErrorResponse struct {
	StatusCode string  `json:"status_code"`
	Error      string  `json:"error"`
	ID         *string `json:"id,omitempty"`
}
