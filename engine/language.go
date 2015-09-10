package engine

import (
	"fmt"
	"io/ioutil"
	"os"

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
}

// RunResult contains the full results for each executed step of a code run.
// Some languages don't have a compile step, and the run step will be nil if
// compilation failed.
type RunResult struct {
	CompileStep, RunStep *ExecutionResult
}

// RunOptions is configuration for a Language Run.
type RunOptions struct {
	Source, Stdin string
	Timeout       int64
	MaxOutputSize int
}

// Run executes the given source code in a sandboxed environment, providing the
// given stdin and returning exit code, stdout and stderr.
func (lang *Language) Run(opts *RunOptions) (*RunResult, error) {
	var exe *execution
	result := &RunResult{}

	dir, err := writeFile(lang.Filename, opts.Source)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	filePath := fmt.Sprintf("/src/%s", lang.Filename)

	if lang.compileStep {
		exe, err = newExecution("compilation", []string{"--build", filePath}, dir, lang.DockerImage, lang.CompilerProfile)
		if err != nil {
			return nil, err
		}
		compileOpts := *opts
		compileOpts.Stdin = ""
		result.CompileStep, err = exe.run(&compileOpts)
		if err != nil || result.CompileStep.ExitCode != 0 {
			return result, err
		}
	}

	exe, err = newExecution("runtime", []string{filePath}, dir, lang.DockerImage, lang.ApparmorProfile)
	if err != nil {
		return nil, err
	}
	result.RunStep, err = exe.run(opts)

	return result, err
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
