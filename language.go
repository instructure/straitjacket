package main

// sudo ./install_apparmor_profiles
// docker build -t cstest languages/csharp
// docker build -t straitjacket . && docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp --rm straitjacket csharp

import(
  "fmt"
  docker "github.com/fsouza/go-dockerclient"
  "code.google.com/p/goconf/conf"
  "os"
  "io/ioutil"
)

type Language struct {
  config *conf.ConfigFile
}

var(
  endpoint = "unix:///var/run/docker.sock"
  tempdir = "/tmp"
)

func LoadLanguage(configName string) (lang Language, err error) {
  confFile := "config/lang-" + configName + ".conf"
  lang.config, err = conf.ReadConfigFile(confFile)
  return
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
  lang, err := LoadLanguage(os.Args[1])
  check(err)

  dir, _ := ioutil.TempDir(tempdir, "straitjacket")
  defer os.RemoveAll(dir)
  check(os.Chmod(dir, 0777))
  filename, _ := lang.config.GetString("general", "filename")
  val, _ := lang.config.GetString("test-simple", "source")
  err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, filename), []byte(val), 0644)
  check(err)

  image, err := lang.config.GetString("general", "docker_image")
  check(err)
  apparmor, err := lang.config.GetString("general", "apparmor_profile")
  check(err)
  client, err := docker.NewClient(endpoint)
  check(err)
  container, err := client.CreateContainer(docker.CreateContainerOptions{
    Config: &docker.Config{
      Image: image,
      Cmd: []string{fmt.Sprintf("/src/%s", filename)},
    },
  })
  check(err)
  defer client.RemoveContainer(docker.RemoveContainerOptions{ ID: container.ID, Force: true })
  err = client.StartContainer(container.ID, &docker.HostConfig{
    Binds: []string{fmt.Sprintf("%s:/src:ro", dir)},
    SecurityOpt: []string{fmt.Sprintf("apparmor:%s", apparmor)},
  })
  check(err)
  status, err := client.WaitContainer(container.ID)
  check(err)
  fmt.Printf("status was: %d\n", status)
  client.Logs(docker.LogsOptions{
    Container: container.ID,
    OutputStream: os.Stdout,
    ErrorStream: os.Stderr,
    Stdout: true,
    Stderr: true,
  })
}