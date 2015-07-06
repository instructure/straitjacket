package main

import(
  "bytes"
  docker "github.com/fsouza/go-dockerclient"
  "io/ioutil"
  "os"
  "fmt"
)

var(
  endpoint = "unix:///var/run/docker.sock"
  tempdir = "/tmp"
)

type RunOptions struct {
  Source, Stdin string
}

type RunResult struct {
  ExitCode int
  Stdout, Stderr string
}

func (lang *Language) Run(opts *RunOptions) (result *RunResult, err error) {
  result = &RunResult{}

  dir, err := ioutil.TempDir(tempdir, "straitjacket")
  if err != nil {
    return
  }
  defer os.RemoveAll(dir)

  check(os.Chmod(dir, 0777))
  err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, lang.filename), []byte(opts.Source), 0644)
  if err != nil {
    return
  }

  client, err := docker.NewClient(endpoint)
  if err != nil {
    return
  }

  container, err := client.CreateContainer(docker.CreateContainerOptions{
    Config: &docker.Config{
      Image: lang.docker_image,
      Cmd: []string{fmt.Sprintf("/src/%s", lang.filename)},
    },
  })
  if err != nil {
    return
  }
  defer client.RemoveContainer(docker.RemoveContainerOptions{ ID: container.ID, Force: true })

  err = client.StartContainer(container.ID, &docker.HostConfig{
    Binds: []string{fmt.Sprintf("%s:/src:ro", dir)},
    SecurityOpt: []string{fmt.Sprintf("apparmor:%s", lang.apparmor_profile)},
  })
  if err != nil {
    return
  }

  result.ExitCode, err = client.WaitContainer(container.ID)
  if err != nil {
    return
  }

  var stdout, stderr bytes.Buffer
  err = client.Logs(docker.LogsOptions{
    Container: container.ID,
    OutputStream: &stdout,
    ErrorStream: &stderr,
    Stdout: true,
    Stderr: true,
  })
  if err != nil {
    return
  }

  result.Stdout = stdout.String()
  result.Stderr = stderr.String()

  return
}
