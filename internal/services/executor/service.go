package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/commander-cli/cmd"
)

type Executioner interface {
	WithEnvironmentVariables(env cmd.EnvVars) Executioner
	WithTimeout(timeout time.Duration) Executioner
	WithWorkingDir(workingDir string) Executioner
	Execute(command string) (result string, err error)
	ExecuteContext(ctx context.Context, command string) (result string, err error)
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
		return "", extractError(base, err)
	}
	return base.Combined(), nil
}

func (e *Executor) ExecuteContext(ctx context.Context, command string) (result string, err error) {
	base := cmd.NewCommand(command, e.options...)
	if err = base.ExecuteContext(ctx); err != nil {
		return "", extractError(base, err)
	}
	return base.Combined(), nil
}

func extractError(command *cmd.Command, err error) error {
	if command.Executed() {
		if stderr := command.Stderr(); stderr != "" {
			err = fmt.Errorf("%w\n\n%s", err, stderr)
		}
	}
	return err
}
