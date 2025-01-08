package set

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/commands/models"
	"github.com/queeck/cli/internal/services/config"
)

const (
	Code = "set"
)

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	bus           commands.Bus
	keymap        keymaps.InputKeymap
	help          help.Model
	inputKey      textinput.Model
	inputValue    textinput.Model
	key           string
	value         string
	quitting      bool
	selected      bool
	done          bool
	isKeyNotFound bool
	isKeyComplex  bool
	setError      string
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

func newInputValue(key string) textinput.Model {
	input := textinput.New()
	input.Placeholder = "<value>"
	input.Prompt = key + " = "
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Focus()
	input.CharLimit = 120
	input.Width = 120
	return input
}

func New(bus commands.Bus) commands.Command {
	return &Model{
		bus:        bus,
		keymap:     keymaps.Input(),
		help:       help.New(),
		inputKey:   newInputKeyWithSuggestions(bus.Config().Keys()),
		inputValue: newInputValue(""),
	}
}

type SelectedKeyMsg struct{}

func SelectedMessage() tea.Cmd {
	return func() tea.Msg {
		return SelectedKeyMsg{}
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
	}
	if index+2 < len(args) {
		m.value = args[index+2]
	}
	if m.key != "" || m.value != "" {
		return SelectedMessage()
	}
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg, tea.MouseMsg:
		return m, nil
	case SelectedKeyMsg:
		return m.choose()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.bus.Parent(m), nil
		case key.Matches(msg, m.keymap.Select):
			return m.enter()
		case key.Matches(msg, m.keymap.Quit):
			return m.quit()
		}
	}

	if !m.selected || m.key == "" {
		m.inputKey, cmd = m.inputKey.Update(msg)
	} else {
		m.inputValue, cmd = m.inputValue.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.setError != "" {
		return withNewLine("Error: " + m.setError)
	}
	if m.quitting {
		return m.bus.Templates().Render(models.TemplateQuit)
	}

	if m.isKeyNotFound {
		return m.bus.Templates().Render(templateKeyNotFound, "key", m.key)
	}

	if m.isKeyComplex {
		return m.bus.Templates().Render(templateKeyHasComplexType, "key", m.key)
	}

	if !m.selected || m.key == "" {
		return m.bus.Templates().Render(templateKeyScreen,
			"inputKey", m.inputKey.View(),
			"help", m.help.View(m.keymap),
		)
	}

	if m.value == "" {
		return m.bus.Templates().Render(templateValueScreen,
			"inputValue", m.inputValue.View(),
			"key", m.key,
			"help", m.help.View(m.keymap),
		)
	}

	return withNewLine(fmt.Sprintf("%s = %s", m.key, m.value))
}

func (m *Model) choose() (tea.Model, tea.Cmd) {
	if m.key != "" {
		m.selected = true
	}
	if m.value != "" {
		m.done = true
	}
	return m.enter()
}

func (m *Model) enter() (tea.Model, tea.Cmd) {
	if !m.selected {
		m.key = m.inputKey.Value()
		m.inputValue.Focus()
		m.selected = true
		return m, nil
	}
	if !m.done {
		m.value = m.inputValue.Value()
	}
	if m.value == "" {
		return m, nil
	}

	valueType := m.bus.Config().Type(m.key)
	if valueType == config.TypeNull {
		m.isKeyNotFound = true
	}
	if valueType == config.TypeObject || valueType == config.TypeArray {
		m.isKeyComplex = true
	}
	if m.isKeyNotFound || m.isKeyComplex {
		return m, tea.Quit
	}
	if err := m.save(); err != nil {
		m.setError = err.Error()
	}
	return m, tea.Quit
}

func (m *Model) save() error {
	if err := m.bus.Config().Set(m.key, m.value); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}
	if err := m.bus.Config().Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
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
