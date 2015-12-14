package engine

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fsouza/go-dockerclient"
	"gopkg.in/yaml.v2"
)

// Language represents a supported execution language, read from the config yml
// files. It provides a Run method for sandboxed code execution.
type Language struct {
	Name            string
	VisibleName     string `yaml:"visible_name"`
	Version         string
	Filename        string
	DockerImage     string `yaml:"docker_image"`
	compileStep     bool
	ApparmorProfile string   `yaml:"apparmor_profile"`
	CompilerProfile string   `yaml:"compiler_profile"`
	Checks          Checks   `yaml:"tests"`
	FileExtensions  []string `yaml:"file_extensions"`
	client          *docker.Client
}

// RunResult contains the full results for each executed step of a code run.
// Some languages don't have a compile step, and the run step will be nil if
// compilation failed.
type RunResult struct {
	CompileStep, RunStep *ExecutionResult
}

// RunOptions is configuration for a Language Run.
type RunOptions struct {
	Source                  string
	Stdin                   io.Reader
	Stdout, Stderr          io.Writer
	Timeout, CompileTimeout int64
	MaxOutputSize           int
}

// Run is a convenience method to compile, execute, and then cleanup code.
func (lang *Language) Run(opts *RunOptions) (result *RunResult, err error) {
	result = &RunResult{}
	image, compileResult, err := lang.Compile(opts.CompileTimeout, opts.Source)
	result.CompileStep = compileResult
	if image != nil {
		defer func() {
			go image.Remove()
		}()

		if compileResult != nil && compileResult.ExitCode != 0 {
			return
		}

		result.RunStep, err = image.Run(opts)
	}

	return
}

func (lang *Language) validate() error {
	if lang.Name == "" {
		return fmt.Errorf("Missing required attribute: name")
	}
	if lang.Filename == "" {
		return fmt.Errorf("Missing required attribute: filename")
	}
	if lang.DockerImage == "" {
		return fmt.Errorf("Missing required attribute: docker_image")
	}
	if lang.ApparmorProfile == "" {
		return fmt.Errorf("Missing required attribute: apparmor_profile")
	}
	return nil
}

func loadLanguage(configName string) (lang *Language, err error) {
	lang = &Language{}
	data, err := ioutil.ReadFile(configName)
	if err == nil {
		err = yaml.Unmarshal(data, lang)
	}
	if lang.CompilerProfile != "" {
		lang.compileStep = true
	}
	if err == nil {
		err = lang.validate()
	}
	return
}

func (lang *Language) disableAppArmor() {
	lang.ApparmorProfile = ""
	lang.CompilerProfile = ""
}
