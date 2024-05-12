package get

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/pkg/styles"
	"github.com/queeck/cli/internal/services/commands"
)

const (
	Code = "get"
)

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	bus      commands.Bus
	keymap   keymaps.InputKeymap
	help     help.Model
	inputKey textinput.Model
	key      string
	value    string
	has      bool
	quitting bool
	selected bool
}

func newInputKey(suggestions []string) textinput.Model {
	input := textinput.New()
	input.Placeholder = "<key.with.dots>"
	input.Prompt = "config/"
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Focus()
	input.SetSuggestions(suggestions)
	input.ShowSuggestions = true
	input.CharLimit = 120
	input.Width = 120
	return input
}

func New(bus commands.Bus) commands.Command {
	return &Model{
		bus:      bus,
		keymap:   keymaps.Input(),
		help:     help.New(),
		inputKey: newInputKey(bus.Config().Keys()),
	}
}

type SelectedMsg struct {
	key string
}

func SelectedMessage(key string) tea.Cmd {
	return func() tea.Msg {
		return SelectedMsg{
			key: key,
		}
	}
}

func (m *Model) Code() string {
	return Code
}

func (m *Model) Commands() []commands.Variant {
	return nil
}

func (m *Model) Init() tea.Cmd {
	args := m.bus.Arguments().Commands()
	index := slices.Index(args, Code)
	if index+1 < len(args) {
		return SelectedMessage(args[index+1])
	}
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg, tea.MouseMsg:
		return m, nil
	case SelectedMsg:
		m.selected = true
		return m.get(msg.key)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.bus.Parent(m), nil
		case key.Matches(msg, m.keymap.Select):
			return m.get(m.inputKey.Value())
		case key.Matches(msg, m.keymap.Quit):
			return m.quit()
		}
	}

	m.inputKey, cmd = m.inputKey.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return m.result() + "\n"
	}

	return fmt.Sprintf(
		"Enter key:\n\n  %s\n\n%s\n\n",
		m.inputKey.View(),
		m.help.View(m.keymap),
	)
}

func (m *Model) result() string {
	if m.selected {
		return m.value
	}
	if !m.has {
		if m.key == "" {
			return styles.ColorForegroundSubtle("(empty key)")
		}
		return styles.ColorForegroundSubtle("(" + m.key + " was not set)")
	}
	return styles.ColorForegroundSubtle(m.key+" = ") + "\n" + m.value + "\n"
}

func (m *Model) get(key string) (tea.Model, tea.Cmd) {
	m.key = key
	m.value, m.has = m.bus.Config().GetString(m.key)
	return m.quit()
}

func (m *Model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}
