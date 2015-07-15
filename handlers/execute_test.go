package handlers

import (
	"straitjacket/engine"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timings = []struct {
	value    string
	expected int64
}{
	{"", 60},
	{"15", 15},
	{"95", 95},
}

func TestParseTimelimit(t *testing.T) {
	for _, param := range timings {
		timelimit, err := parseTimelimit(param.value)
		if assert.NoError(t, err) {
			assert.Equal(t, param.expected, timelimit)
		}
	}
}

func TestFullSuccess(t *testing.T) {
	runResult := &engine.RunResult{
		CompileStep: &engine.ExecutionResult{
			ExitCode:    0,
			RunTime:     1 * time.Second,
			ErrorString: "",
			Stdout:      "x",
			Stderr:      "y",
		},
		RunStep: &engine.ExecutionResult{
			ExitCode:    0,
			RunTime:     time.Duration(2.5 * float64(time.Second)),
			ErrorString: "",
			Stdout:      "rx",
			Stderr:      "ry",
		},
	}
	result := buildResult(runResult)
	assert.Equal(t, true, result.Success)
	assert.Nil(t, result.Error)
	assert.Nil(t, result.Compilation.Error)
	assert.Nil(t, result.Runtime.Error)
	assert.Equal(t, 0, result.Compilation.ExitStatus)
	assert.Equal(t, 0, result.Runtime.ExitStatus)
	assert.Equal(t, 1.0, result.Compilation.Time)
	assert.Equal(t, 2.5, result.Runtime.Time)
	assert.Equal(t, "x", result.Compilation.Stdout)
	assert.Equal(t, "y", result.Compilation.Stderr)
	assert.Equal(t, "rx", result.Runtime.Stdout)
	assert.Equal(t, "ry", result.Runtime.Stderr)
}

func TestNoCompileStep(t *testing.T) {
	runResult := &engine.RunResult{
		RunStep: &engine.ExecutionResult{
			ExitCode:    0,
			RunTime:     time.Duration(2.5 * float64(time.Second)),
			ErrorString: "",
			Stdout:      "x",
			Stderr:      "y",
		},
	}
	result := buildResult(runResult)
	assert.Equal(t, true, result.Success)
	assert.Nil(t, result.Error)
	assert.Nil(t, result.Compilation)
	assert.Nil(t, result.Runtime.Error)
	assert.Equal(t, 0, result.Runtime.ExitStatus)
	assert.Equal(t, 2.5, result.Runtime.Time)
	assert.Equal(t, "x", result.Runtime.Stdout)
	assert.Equal(t, "y", result.Runtime.Stderr)
}

func TestFailedCompileStep(t *testing.T) {
	runResult := &engine.RunResult{
		CompileStep: &engine.ExecutionResult{
			ExitCode:    3,
			RunTime:     time.Duration(2.5 * float64(time.Second)),
			ErrorString: "compilation_error",
			Stdout:      "x",
			Stderr:      "y",
		},
	}
	result := buildResult(runResult)
	assert.Equal(t, false, result.Success)
	assert.Equal(t, "compilation_error", *result.Error)
	assert.Nil(t, result.Runtime)
	assert.Equal(t, "compilation_error", *result.Compilation.Error)
	assert.Equal(t, 3, result.Compilation.ExitStatus)
	assert.Equal(t, 2.5, result.Compilation.Time)
	assert.Equal(t, "x", result.Compilation.Stdout)
	assert.Equal(t, "y", result.Compilation.Stderr)
}

func TestFailedRuntimeStep(t *testing.T) {
	runResult := &engine.RunResult{
		CompileStep: &engine.ExecutionResult{
			ExitCode:    0,
			RunTime:     time.Duration(2.5 * float64(time.Second)),
			ErrorString: "",
			Stdout:      "x",
			Stderr:      "y",
		},
		RunStep: &engine.ExecutionResult{
			ExitCode:    3,
			RunTime:     1500,
			ErrorString: "runtime_error",
			Stdout:      "rx",
			Stderr:      "ry",
		},
	}
	result := buildResult(runResult)
	assert.Equal(t, false, result.Success)
	assert.Nil(t, result.Compilation.Error)
	assert.Equal(t, "runtime_error", *result.Error)
	assert.Equal(t, 0, result.Compilation.ExitStatus)
	assert.Equal(t, 3, result.Runtime.ExitStatus)
}
