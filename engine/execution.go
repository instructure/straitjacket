package engine

import (
	"fmt"
	"io"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	endpoint = "unix:///var/run/docker.sock"
)

// ExecutionResult contains the results from one step of code execution --
// either compilation or running.
type ExecutionResult struct {
	// The exit status code of the process. 0 is success, anything else is failure.
	ExitCode int
	// How long this step ran.
	RunTime time.Duration
	// If an error occured, this error string will be non-empty.
	ErrorString    string
	Stdout, Stderr string
}

type executionOptions struct {
	apparmorProfile string
	stdin           io.Reader
	stdout, stderr  io.Writer
	timeout         int64
	maxOutputSize   int
}

type container struct {
	id     string
	client *docker.Client
}

func createContainer(client *docker.Client, opts docker.CreateContainerOptions) (c *container, err error) {
	var dockerInfo *docker.Container
	dockerInfo, err = client.CreateContainer(opts)
	if err == nil {
		c = &container{
			client: client,
			id:     dockerInfo.ID,
		}
	}
	return
}

func (c *container) Remove() {
	c.client.RemoveContainer(docker.RemoveContainerOptions{ID: c.id, Force: true})
}

func (c *container) execute(step string, opts *executionOptions) (result *ExecutionResult, err error) {
	result = &ExecutionResult{}
	timeout := false

	if err == nil {
		startTime := time.Now()
		err1, err2 := c.attachAndRun(result, opts)
		select {
		case err = <-err1:
		case err = <-err2:
			// pass
		case <-time.After(time.Duration(opts.timeout) * time.Second):
			timeout = true
		}
		result.RunTime = time.Now().Sub(startTime)

		// if the container is already stopped we just ignore the error respose...
		// RemoveContainer will be called later.
		c.client.KillContainer(docker.KillContainerOptions{ID: c.id})
		// make sure we've shut down by waiting for attachAndRun to push or close the channels
		err = firstError(err, <-err1, <-err2)
	}

	if timeout {
		result.ErrorString = fmt.Sprintf("%s_timelimit", step)
		result.ExitCode = -9
	} else if _, ok := err.(*OutputTooLarge); ok {
		// treat too-large output as a soft error, still returning a response
		err = nil
		result.ExitCode = -10
		result.ErrorString = fmt.Sprintf("%s_output_size_error", step)
	} else if result.ExitCode != 0 {
		result.ErrorString = fmt.Sprintf("%s_error", step)
	}

	return
}

func (c *container) attachAndRun(result *ExecutionResult, opts *executionOptions) (chan error, chan error) {
	sentinel := make(chan struct{})
	runResult := make(chan error)
	attachResult := make(chan error)

	// run goroutine
	go func() {
		_ = <-sentinel
		securityOpt := []string{}
		if opts.apparmorProfile != "" {
			securityOpt = append(securityOpt, fmt.Sprintf("apparmor:%s", opts.apparmorProfile))
		}
		// when we get the sentinel, we know we've attached in the other goroutine
		err := c.client.StartContainer(c.id, &docker.HostConfig{
			SecurityOpt: securityOpt,
			LogConfig: docker.LogConfig{
				Type: "none",
			},
		})
		sentinel <- struct{}{}

		if err == nil {
			result.ExitCode, err = c.client.WaitContainer(c.id)
		}
		runResult <- err
		close(runResult)
	}()

	// attach goroutine
	go func() {
		err := c.client.AttachToContainer(docker.AttachToContainerOptions{
			Container:    c.id,
			InputStream:  opts.stdin,
			OutputStream: NewLimitedWriter(opts.stdout, opts.maxOutputSize),
			ErrorStream:  NewLimitedWriter(opts.stderr, opts.maxOutputSize),
			Stream:       true,
			Stdin:        true,
			Stdout:       true,
			Stderr:       true,
			Success:      sentinel,
		})

		attachResult <- err
		close(attachResult)
	}()

	return runResult, attachResult
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
