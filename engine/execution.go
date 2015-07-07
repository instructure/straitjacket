package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	endpoint = "unix:///var/run/docker.sock"
	tempdir  = "/tmp"
)

type RunOptions struct {
	Source, Stdin string
}

type RunResult struct {
	ExitCode       int
	Stdout, Stderr string
}

func (lang *Language) Run(opts *RunOptions) (result *RunResult, err error) {
	result = &RunResult{}

	dir, err := ioutil.TempDir(tempdir, "straitjacket")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	err = os.Chmod(dir, 0777)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, lang.Filename), []byte(opts.Source), 0644)
	if err != nil {
		return
	}

	client, err := docker.NewClient(endpoint)
	if err != nil {
		return
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:     lang.DockerImage,
			Cmd:       []string{fmt.Sprintf("/src/%s", lang.Filename)},
			OpenStdin: true,
			StdinOnce: true,
		},
	})
	if err != nil {
		return
	}
	defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID, Force: true})

	sentinel := make(chan struct{})
	runResult := make(chan error)

	go func() {
		_ = <-sentinel
		// when we get the sentinel, we know we've attached in the main thread
		err := client.StartContainer(container.ID, &docker.HostConfig{
			Binds:       []string{fmt.Sprintf("%s:/src:ro", dir)},
			SecurityOpt: []string{fmt.Sprintf("apparmor:%s", lang.ApparmorProfile)},
		})
		sentinel <- struct{}{}

		if err == nil {
			result.ExitCode, err = client.WaitContainer(container.ID)
		}
		runResult <- err
	}()

	stdin := strings.NewReader(opts.Stdin)
	var stdout, stderr bytes.Buffer
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		InputStream:  stdin,
		OutputStream: &stdout,
		ErrorStream:  &stderr,
		Stream:       true,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		Success:      sentinel,
	})

	if err == nil {
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
		err = <-runResult
	}

	return
}
