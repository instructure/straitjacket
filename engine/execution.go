package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	endpoint = "unix:///var/run/docker.sock"
	tempdir  = "/tmp"
)

type RunOptions struct {
	Source, Stdin string
	Timeout       int64
}

type RunResult struct {
	ExitCode       int
	Stdout, Stderr string
	RunTime        time.Duration
	ErrorString    string
}

type execution struct {
	lang      *Language
	tmpDir    string
	client    *docker.Client
	container *docker.Container
	sentinel  chan struct{}
	result    *RunResult
}

// Initialize a new exeuction object for use.
func newExecution(lang *Language) (exe *execution, err error) {
	exe = &execution{
		lang:     lang,
		sentinel: make(chan struct{}),
		result:   &RunResult{},
	}

	exe.client, err = docker.NewClient(endpoint)
	return
}

// Run the execution with the given options.
func (exe *execution) run(opts *RunOptions) (result *RunResult, err error) {
	result = exe.result
	timeout := false
	defer exe.cleanup()

	exe.tmpDir, err = writeFile(exe.lang.Filename, opts.Source)

	if err == nil {
		err = exe.createContainer()
	}

	if err == nil {
		startTime := time.Now()
		runResult := exe.attachAndRun(opts.Stdin)
		select {
		case err = <-runResult:
			// pass
		case <-time.After(time.Duration(opts.Timeout) * time.Second):
			timeout = true
		}
		result.RunTime = time.Now().Sub(startTime)
	}

	if timeout {
		result.ErrorString = "runtime_timelimit"
	} else if result.ExitCode != 0 {
		result.ErrorString = "runtime_error"
	}

	return
}

func (exe *execution) createContainer() (err error) {
	exe.container, err = exe.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:     exe.lang.DockerImage,
			Cmd:       []string{fmt.Sprintf("/src/%s", exe.lang.Filename)},
			OpenStdin: true,
			StdinOnce: true,
		},
	})

	return
}

func (exe *execution) attachAndRun(stdin string) chan error {
	sentinel := make(chan struct{})
	runResult := make(chan error)
	attachResult := make(chan error)
	finalResult := make(chan error)

	// run goroutine
	go func() {
		_ = <-sentinel
		// when we get the sentinel, we know we've attached in the other goroutine
		err := exe.client.StartContainer(exe.container.ID, &docker.HostConfig{
			Binds:       []string{fmt.Sprintf("%s:/src:ro", exe.tmpDir)},
			SecurityOpt: []string{fmt.Sprintf("apparmor:%s", exe.lang.ApparmorProfile)},
		})
		sentinel <- struct{}{}

		if err == nil {
			exe.result.ExitCode, err = exe.client.WaitContainer(exe.container.ID)
		}
		runResult <- err
	}()

	// attach goroutine
	go func() {
		stdinReader := strings.NewReader(stdin)
		var stdout, stderr bytes.Buffer
		err := exe.client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    exe.container.ID,
			InputStream:  stdinReader,
			OutputStream: &stdout,
			ErrorStream:  &stderr,
			Stream:       true,
			Stdin:        true,
			Stdout:       true,
			Stderr:       true,
			Success:      sentinel,
		})

		if err == nil {
			exe.result.Stdout = stdout.String()
			exe.result.Stderr = stderr.String()
		}

		attachResult <- err
	}()

	go func() {
		finalResult <- firstError(<-runResult, <-attachResult)
	}()

	return finalResult
}

func writeFile(filename, source string) (string, error) {
	dir, err := ioutil.TempDir(tempdir, "straitjacket")

	if err == nil {
		err = os.Chmod(dir, 0777)
	}
	if err == nil {
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, filename), []byte(source), 0644)
	}

	if err != nil {
		dir = ""
	}

	return dir, err
}

func (exe *execution) cleanup() {
	if exe.tmpDir != "" {
		os.RemoveAll(exe.tmpDir)
	}
	if exe.container != nil {
		exe.client.RemoveContainer(docker.RemoveContainerOptions{ID: exe.container.ID, Force: true})
	}
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
