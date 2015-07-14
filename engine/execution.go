package engine

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	endpoint = "unix:///var/run/docker.sock"
	tempdir  = "/tmp"
)

type ExecutionResult struct {
	ExitCode       int
	Stdout, Stderr string
	RunTime        time.Duration
	ErrorString    string
}

type execution struct {
	step            string
	command         []string
	srcDir          string
	dockerImage     string
	apparmorProfile string
	client          *docker.Client
	container       *docker.Container
	sentinel        chan struct{}
	result          *ExecutionResult
}

// Initialize a new exeuction object for use.
func newExecution(step string, command []string, srcDir, dockerImage, apparmorProfile string) (exe *execution, err error) {
	exe = &execution{
		step:            step,
		command:         command,
		srcDir:          srcDir,
		dockerImage:     dockerImage,
		apparmorProfile: apparmorProfile,
		sentinel:        make(chan struct{}),
		result:          &ExecutionResult{},
	}

	exe.client, err = docker.NewClient(endpoint)
	return
}

// Run the execution with the given options.
func (exe *execution) run(opts *RunOptions) (result *ExecutionResult, err error) {
	result = exe.result
	timeout := false
	defer exe.cleanup()

	err = exe.createContainer()

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
		result.ErrorString = fmt.Sprintf("%s_timelimit", exe.step)
		result.ExitCode = -9
	} else if result.ExitCode != 0 {
		result.ErrorString = fmt.Sprintf("%s_error", exe.step)
	}

	return
}

func (exe *execution) createContainer() (err error) {
	exe.container, err = exe.client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:     exe.dockerImage,
			Cmd:       exe.command,
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
		securityOpt := []string{}
		if exe.apparmorProfile != "" {
			securityOpt = append(securityOpt, fmt.Sprintf("apparmor:%s", exe.apparmorProfile))
		}
		// when we get the sentinel, we know we've attached in the other goroutine
		err := exe.client.StartContainer(exe.container.ID, &docker.HostConfig{
			Binds:       []string{fmt.Sprintf("%s:/src", exe.srcDir)},
			SecurityOpt: securityOpt,
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

func (exe *execution) cleanup() {
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
