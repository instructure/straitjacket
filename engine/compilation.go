package engine

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"

	docker "github.com/fsouza/go-dockerclient"
)

// An Image is the result of a compile, a reference to a docker container
// containing the uploaded and compiled source code. It can be used to spawn
// multiple execution containers (internally this works via a shared read-only
// volume).
//
// Originally we did this by commiting a new docker image containing the
// compiled source, but that docker commit step takes over 1s on average, so now
// we use the shared volume approach.
type Image struct {
	ID     string
	lang   *Language
	client *docker.Client
}

func (image *Image) Run(opts *RunOptions) (result *ExecutionResult, err error) {
	filePath := fmt.Sprintf("/src/%s", image.lang.Filename)
	container, err := createContainer(image.client, docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           image.lang.DockerImage,
			Cmd:             []string{filePath},
			OpenStdin:       true,
			StdinOnce:       true,
			NetworkDisabled: true,
		},
		HostConfig: &docker.HostConfig{
			VolumesFrom: []string{image.ID + ":ro"},
		},
	})

	if err == nil {
		defer func() {
			go container.Remove()
		}()
		result = &ExecutionResult{}
		result, err = container.execute("runtime", &executionOptions{
			timeout:         opts.Timeout,
			stdin:           opts.Stdin,
			stdout:          opts.Stdout,
			stderr:          opts.Stderr,
			maxOutputSize:   opts.MaxOutputSize,
			apparmorProfile: image.lang.ApparmorProfile,
		})
	}

	return
}

// Remove will delete this image from docker, make sure this gets called so we
// don't end up keeping images around.
func (image *Image) Remove() {
	image.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            image.ID,
		RemoveVolumes: true,
		Force:         true,
	})
}

// Compile sets up a new Container and gets ready to execute the given source code.
// It's important to Remove the Container to clean up resources.
func (lang *Language) Compile(timeout int64, source string) (image *Image, result *ExecutionResult, err error) {
	filePath := fmt.Sprintf("/src/%s", lang.Filename)

	container, err := createContainer(lang.client, docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           lang.DockerImage,
			Cmd:             []string{"--build", filePath},
			OpenStdin:       true,
			StdinOnce:       true,
			NetworkDisabled: true,
			Volumes: map[string]struct{}{
				"/src": struct{}{},
			},
		},
	})

	if err == nil {
		err = lang.client.UploadToContainer(container.id, docker.UploadToContainerOptions{
			InputStream: lang.tarSource(source, filePath),
			Path:        "/",
		})
	}

	result = &ExecutionResult{}

	if err == nil && lang.compileStep {
		var stdout, stderr bytes.Buffer
		result, err = container.execute("compilation", &executionOptions{
			timeout:         timeout,
			stdout:          &stdout,
			stderr:          &stderr,
			maxOutputSize:   64 * 1024,
			apparmorProfile: lang.CompilerProfile,
		})
		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
	}

	if err == nil && result.ExitCode == 0 {
		image = &Image{
			ID:     container.id,
			lang:   lang,
			client: lang.client,
		}
	}

	if container != nil && image == nil {
		lang.client.RemoveContainer(docker.RemoveContainerOptions{ID: container.id, Force: true})
	}

	return
}

func (lang *Language) tarSource(source, filePath string) io.Reader {
	result := &bytes.Buffer{}
	writer := tar.NewWriter(result)
	writer.WriteHeader(&tar.Header{
		Name:     "/src",
		Mode:     0777,
		Typeflag: tar.TypeDir,
	})
	writer.WriteHeader(&tar.Header{
		Name: filePath,
		Mode: 0444,
		Size: int64(len(source)),
	})
	writer.Write([]byte(source))
	writer.Close()
	return result
}
