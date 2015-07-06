package handlers

import (
	"net/http"
)

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
										<option value="javascript">JavaScript</option>
										<option value="csharp">C#</option>
										<option value="ruby">Ruby</option>
								</select>
						</dd>
				</dl>
				<input type="submit"/>
		</form>
  </body>
</html>
`

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	_, err := res.Write([]byte(WELCOME_HTML))
	if err != nil {
		panic(err)
	}
}
