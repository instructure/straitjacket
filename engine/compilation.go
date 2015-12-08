package engine

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"

	docker "github.com/fsouza/go-dockerclient"
)

type Image struct {
	ID     string
	lang   *Language
	client *docker.Client
}

func (image *Image) Run(opts *RunOptions) (result *ExecutionResult, err error) {
	filePath := fmt.Sprintf("/src/%s", image.lang.Filename)
	container, err := createContainer(image.client, docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           image.ID,
			Cmd:             []string{filePath},
			OpenStdin:       true,
			StdinOnce:       true,
			NetworkDisabled: true,
		},
	})

	if err == nil {
		defer container.Remove()
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

func (image *Image) Remove() {
	image.client.RemoveImageExtended(image.ID, docker.RemoveImageOptions{
		Force: true,
	})
}

// Compile sets up a new image ready to execute the given source code.
// It's important to Remove the image to clean up resources.
func (lang *Language) Compile(timeout int64, source string) (image *Image, result *ExecutionResult, err error) {
	client, _ := docker.NewClient(endpoint)

	filePath := fmt.Sprintf("/src/%s", lang.Filename)

	container, err := createContainer(client, docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:           lang.DockerImage,
			Cmd:             []string{"--build", filePath},
			OpenStdin:       true,
			StdinOnce:       true,
			NetworkDisabled: true,
		},
	})

	if err == nil {
		defer client.RemoveContainer(docker.RemoveContainerOptions{ID: container.id, Force: true})

		err = client.UploadToContainer(container.id, docker.UploadToContainerOptions{
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
		var dockerImage *docker.Image
		dockerImage, err = client.CommitContainer(docker.CommitContainerOptions{
			Container: container.id,
		})
		if err == nil {
			image = &Image{
				ID:     dockerImage.ID,
				lang:   lang,
				client: client,
			}
		}
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
