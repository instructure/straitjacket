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
    if lang.Name == os.Args[1] {
      test := lang.Tests.Simple
      result, err := lang.Run(&RunOptions{
        Source: test.Source,
        Stdin: test.Stdin,
      })
      check(err)
      fmt.Printf("status was: %d\n", result.ExitCode)
      fmt.Printf("%s", result.Stdout)
      fmt.Printf("%s", result.Stderr)
      return
    }
  }

  fmt.Printf("language %s not found\n", os.Args[1])
}