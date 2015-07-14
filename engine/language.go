package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

type test struct {
	Source, Stdin, Stdout, Stderr string
	ExitStatus                    int
}

type tests struct {
	Simple, Apparmor, Rlimit test
}

// Language represents a supported execution language, read from the config yml
// files. It provides a Run method for sandboxed code execution.
type Language struct {
	Name            string
	VisibleName     string `yaml:"visible_name"`
	Version         string
	Filename        string
	DockerImage     string `yaml:"docker_image"`
	compileStep     bool
	ApparmorProfile string `yaml:"apparmor_profile"`
	CompilerProfile string `yaml:"compiler_profile"`
	Tests           tests
	FileExtensions  []string `yaml:"file_extensions"`
}

type RunResult struct {
	ExitCode       int
	Stdout, Stderr string
	RunTime        time.Duration
	ErrorString    string
}

type RunOptions struct {
	Source, Stdin string
	Timeout       int64
}

// Run executes the given source code in a sandboxed environment, providing the
// given stdin and returning exit code, stdout and stderr.
func (lang *Language) Run(opts *RunOptions) (*RunResult, error) {
	var exe *execution
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
		result, err := exe.run(opts)
		if err != nil || result.ExitCode != 0 {
			return result, err
		}
	}

	exe, err = newExecution("runtime", []string{filePath}, dir, lang.DockerImage, lang.ApparmorProfile)
	if err != nil {
		return nil, err
	}
	result, err := exe.run(opts)

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

// RunTests does sanity checks on all supported Languages, including checks that
// stdin/stdout work as expected, and basic verification that the AppArmor
// profile is in effect.
func (lang *Language) RunTests() (err error) {
	err = lang.runTest("simple", &lang.Tests.Simple)
	if lang.ApparmorProfile == "" {
		// skip the apparmor tests when it's been disabled
		return
	}

	if err == nil {
		err = lang.runTest("apparmor", &lang.Tests.Apparmor)
	}
	if err == nil {
		err = lang.runTest("rlimit", &lang.Tests.Rlimit)
	}
	return
}

func (lang *Language) runTest(testName string, test *test) error {
	result, err := lang.Run(&RunOptions{
		Source:  test.Source,
		Stdin:   test.Stdin,
		Timeout: 30,
	})

	if err != nil {
		return fmt.Errorf("Failure testing '%s' (%s): %v", lang.Name, testName, err)
	}

	errorString := fmt.Sprintf("for '%s' (%s).\nexit code: %d\nstdout: %v\nstderr: %v\nerror: %s\n", lang.Name, testName, result.ExitCode, result.Stdout, result.Stderr, result.ErrorString)

	if result.ExitCode != test.ExitStatus {
		return fmt.Errorf("Incorrect exit code %s", errorString)
	}

	match, err := regexp.MatchString(test.Stderr, result.Stderr)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Incorrect stderr %s", errorString)
	}

	match, err = regexp.MatchString(test.Stdout, result.Stdout)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Incorrect stdout %s", errorString)
	}

	return nil
}

func (lang *Language) disableAppArmor() {
	lang.ApparmorProfile = ""
	lang.CompilerProfile = ""
}
