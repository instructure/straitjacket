package engine

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Test struct {
	Source, Stdin, Stdout, Stderr string
	ExitStatus                    int
}

type Tests struct {
	Simple, Apparmor, Rlimit Test
}

type Language struct {
	Name            string
	VisibleName     string `yaml:"visible_name"`
	Version         string
	Filename        string
	DockerImage     string `yaml:"docker_image"`
	ApparmorProfile string `yaml:"apparmor_profile"`
	Tests           Tests
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

func LoadLanguage(configName string) (lang *Language, err error) {
	lang = &Language{}
	data, err := ioutil.ReadFile(configName)
	if err == nil {
		err = yaml.Unmarshal(data, lang)
	}
	if err == nil {
		err = lang.validate()
	}
	return
}

func (lang *Language) runTests() (err error) {
	err = lang.runTest("simple", &lang.Tests.Simple)
	if err == nil {
		err = lang.runTest("apparmor", &lang.Tests.Apparmor)
	}
	if err == nil {
		err = lang.runTest("rlimit", &lang.Tests.Rlimit)
	}
	return
}

func (lang *Language) runTest(testName string, test *Test) error {
	result, err := lang.Run(&RunOptions{
		Source: test.Source,
		Stdin:  test.Stdin,
	})

	if err != nil {
		return err
	}

	if result.ExitCode != test.ExitStatus {
		return fmt.Errorf("Failure testing '%s' (%s), expected exit status: %d got: %d", lang.Name, testName, test.ExitStatus, result.ExitCode)
	}

	match, err := regexp.MatchString(test.Stderr, result.Stderr)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Failure testing '%s' (%s), got stderr: %o", lang.Name, testName, result.Stderr)
	}

	match, err = regexp.MatchString(test.Stdout, result.Stdout)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Failure testing '%s' (%s), got stdout: %o", lang.Name, testName, result.Stdout)
	}

	return nil
}
