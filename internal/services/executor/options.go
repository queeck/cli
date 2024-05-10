package executor

import (
	"time"

	"github.com/commander-cli/cmd"
)

type Option func(executor *Executor)

func WithEnvironmentVariables(env cmd.EnvVars) Option {
	return func(e *Executor) {
		e.options = append(e.options, cmd.WithEnvironmentVariables(env))
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(e *Executor) {
		e.options = append(e.options, cmd.WithTimeout(timeout))
	}
}

func WithWorkingDir(workingDir string) Option {
	return func(e *Executor) {
		e.options = append(e.options, cmd.WithWorkingDir(workingDir))
	}
}
