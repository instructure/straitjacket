package handlers

import (
	"html/template"
	"net/http"
	"straitjacket/engine"
)

type options struct {
	Languages []*engine.Language
}

const WELCOME_HTML = `
<html>
  <head>
    <title>Straitjacket</title>
  </head>
  <body>
    <h1>Welcome to Straitjacket</h1>
    <form method="post" action="/execute">
				<dl>
						<dt>source:</dt>
						<dd><textarea name="source" rows=8 cols=50></textarea></dd>
						<dt>stdin:</dt>
						<dd><textarea name="stdin" rows=8 cols=50></textarea></dd>
						<dt>language:</dt>
						<dd>
								<select name="language">
									{{range .Languages}}
											<option value="{{.Name}}">{{.VisibleName}}</option>
									{{end}}
								</select>
						</dd>
				</dl>
				<input type="submit"/>
		</form>
  </body>
</html>
`

func (ctx *Context) IndexHandler(res http.ResponseWriter, req *http.Request) {
	opts := &options{
		Languages: ctx.Engine.Languages(),
	}
	tmpl, err := template.New("index").Parse(WELCOME_HTML)
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(res, opts)
	if err != nil {
		panic(err)
	}
}
