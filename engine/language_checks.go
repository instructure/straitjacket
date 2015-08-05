package engine

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

// RunChecks does sanity checks on the Language, including checks that
// stdin/stdout work as expected, and basic verification that the AppArmor
// profile is in effect.
func (lang *Language) RunChecks() (err error) {
	err = lang.runCheck("simple", &lang.Checks.Simple)
	if lang.ApparmorProfile == "" {
		// skip the apparmor tests when it's been disabled
		return
	}

	if err == nil {
		err = lang.runCheck("apparmor", &lang.Checks.Apparmor)
	}
	if err == nil {
		err = lang.runCheck("rlimit", &lang.Checks.Rlimit)
	}
	return
}

type Check struct {
	Source, Stdin, Stdout, Stderr string
	ExitStatus                    int
}

type Checks struct {
	Simple, Apparmor, Rlimit Check
}

func (lang *Language) runCheck(testName string, check *Check) error {
	result, err := lang.Run(&RunOptions{
		Source:  check.Source,
		Stdin:   check.Stdin,
		Timeout: 30,
	})

	if err != nil {
		return fmt.Errorf("Failure testing '%s' (%s): %v", lang.Name, testName, err)
	}

	output, _ := yaml.Marshal(result)

	errorString := fmt.Sprintf("for '%s' (%s).\n%s", lang.Name, testName, output)

	if result.RunStep == nil {
		return fmt.Errorf("Didn't run %s", errorString)
	}

	if result.RunStep.ExitCode != check.ExitStatus {
		return fmt.Errorf("Incorrect exit code %s", errorString)
	}

	match, err := regexp.MatchString(check.Stderr, result.RunStep.Stderr)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Incorrect stderr %s", errorString)
	}

	match, err = regexp.MatchString(check.Stdout, result.RunStep.Stdout)
	if err != nil {
		return err
	}
	if match == false {
		return fmt.Errorf("Incorrect stdout %s", errorString)
	}

	return nil
}
