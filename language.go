package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Test struct {
	Source, Stdin, Stdout, Stderr, ExitStatus string
}

type Tests struct {
	Simple, Apparmor, Rlimit Test
}

type Language struct {
	Name            string
	VisibleName     string `yaml:"visible_name"`
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
