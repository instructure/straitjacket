package main

import "straitjacket/engine"

// sudo ./install_apparmor_profiles
// docker build -t cstest languages/csharp
// docker build -t straitjacket . && docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp --rm -p 8081:8081 straitjacket

func main() {
	var err error
	engine.TheEngine, err = engine.LoadConfig("config")
	if err != nil {
		panic(err)
	}
	StartServer(":8081")
}
