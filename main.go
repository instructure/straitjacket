package main

// sudo ./install_apparmor_profiles
// docker build -t cstest languages/csharp
// docker build -t straitjacket . && docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp --rm straitjacket C#

import(
  "fmt"
  "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
  straitjacket, err := LoadConfig("config")
  check(err)
  for _, lang := range straitjacket.Languages {
    if lang.name == os.Args[1] {
      test_source, _ := lang.config.GetString("test-simple", "source")
      result, err := lang.Run(&RunOptions{
        Source: test_source,
      })
      check(err)
      fmt.Printf("status was: %d\n", result.ExitCode)
      fmt.Printf("%s", result.Stdout)
      return
    }
  }

  fmt.Printf("language %s not found", os.Args[1])
}