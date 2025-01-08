package get

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/keymaps"
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

func newInputKeyWithSuggestions(suggestions []string) textinput.Model {
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
		inputKey: newInputKeyWithSuggestions(bus.Config().Keys()),
	}
}

type SelectedMsg struct{}

func SelectedMessage() tea.Cmd {
	return func() tea.Msg {
		return SelectedMsg{}
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
		m.key = args[index+1]
		return SelectedMessage()
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
		return m.get()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.bus.Parent(m), nil
		case key.Matches(msg, m.keymap.Select):
			m.key = m.inputKey.Value()
			return m.get()
		case key.Matches(msg, m.keymap.Quit):
			return m.quit()
		}
	}

	m.inputKey, cmd = m.inputKey.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	if m.quitting {
		return withNewLine(m.result())
	}

	return m.bus.Templates().Render(templateScreen,
		"inputKey", m.inputKey.View(),
		"help", m.help.View(m.keymap),
	)
}

func (m *Model) result() string {
	if m.selected {
		return m.bus.Templates().Render(templateValue, "value", m.value)
	}
	if !m.has {
		if m.key == "" {
			return m.bus.Templates().Render(templateKeyIsEmpty)
		}
		return m.bus.Templates().Render(templateKeyNotFound, "key", m.key)
	}
	return m.bus.Templates().Render(templateValueWithKey, "key", m.key, "value", m.value)
}

func (m *Model) get() (tea.Model, tea.Cmd) {
	m.value, m.has = m.bus.Config().GetString(m.key)
	return m.quit()
}

func (m *Model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func withNewLine(text string) string {
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text
}
