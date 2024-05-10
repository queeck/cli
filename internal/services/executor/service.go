package executor

import (
	"fmt"
	"time"

	"github.com/commander-cli/cmd"
)

type Executioner interface {
	WithEnvironmentVariables(env cmd.EnvVars) Executioner
	WithTimeout(timeout time.Duration) Executioner
	WithWorkingDir(workingDir string) Executioner
	Execute(command string) (result string, err error)
}

type Executor struct {
	options []func(c *cmd.Command)
}

func New(options ...Option) *Executor {
	e := &Executor{
		options: []func(c *cmd.Command){
			cmd.WithInheritedEnvironment(nil),
		},
	}
	for _, option := range options {
		option(e)
	}
	return e
}

func (e *Executor) WithEnvironmentVariables(env cmd.EnvVars) Executioner {
	n := *e
	n.options = append(n.options, cmd.WithEnvironmentVariables(env))
	return &n
}

func (e *Executor) WithTimeout(timeout time.Duration) Executioner {
	n := *e
	n.options = append(n.options, cmd.WithTimeout(timeout))
	return &n
}

func (e *Executor) WithWorkingDir(workingDir string) Executioner {
	n := *e
	n.options = append(n.options, cmd.WithWorkingDir(workingDir))
	return &n
}

func (e *Executor) Execute(command string) (result string, err error) {
	base := cmd.NewCommand(command, e.options...)
	if err = base.Execute(); err != nil {
		if base.Executed() {
			if stderr := base.Stderr(); stderr != "" {
				err = fmt.Errorf("%w\n\n%s", err, stderr)
			}
		}
		return "", err
	}
	return base.Combined(), nil
}

func Execute(command string, options ...Option) (result string, err error) {
	return New(options...).Execute(command)
}

func MustExecute(command string, options ...Option) string {
	result, err := New(options...).Execute(command)
	if err != nil {
		panic(err)
	}
	return result
}
