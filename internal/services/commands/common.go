package commands

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/queeck/cli/internal/pkg/cli"
	"github.com/queeck/cli/internal/services/config"
	"github.com/queeck/cli/internal/services/state"
	"github.com/queeck/cli/internal/services/templates"
)

type Variant struct {
	Code        string
	Description string
}

type Command interface {
	tea.Model
	Code() string
	Commands() []Variant
}

// Bus is common interface for all commands.
// Here is for prevent import cycle, this interface uses by:
// - services/stack
// - services/commands/...
type Bus interface {
	State() state.State
	Arguments() cli.Arguments
	Config() config.Config
	Templates() templates.Renderer
	Parent(command Command) Command
	Children(command Command) []Command
	Child(command Command, code string) Command
	CommandByCLIArguments(arguments cli.Arguments) Command
	SelectedCommands(command Command, code string) string
	CommandRoot() Command
	CommandConfig() Command
	CommandConfigView() Command
}
