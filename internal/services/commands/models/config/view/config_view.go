package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/models"
	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/services/commands"
)

const (
	Code = "view"
)

var _ commands.Command = &ConfigView{} // check for interface compatibility

type ConfigView struct {
	bus      commands.Bus
	keys     keymaps.ViewportKeyMap
	help     help.Model
	viewport viewport.Model
	view     string
	quitting bool
	ready    bool
	width    int
	height   int
}

func New(bus commands.Bus) commands.Command {
	return &ConfigView{
		view:  bus.Config().View(),
		keys:  keymaps.Viewport(),
		help:  help.New(),
		bus:   bus,
		ready: false,
	}
}

func (m *ConfigView) Code() string {
	return Code
}

func (m *ConfigView) Commands() []models.Command {
	return nil
}

func (m *ConfigView) Init() tea.Cmd {
	return nil
}

func (m *ConfigView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = m.width
		m.sync()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			return m.bus.CommandConfig(), nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.ready = false
			m.sync()
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, tea.Batch(cmd)
}

func (m *ConfigView) View() string {
	if m.quitting {
		return m.bus.Templates().RenderCommonQuit()
	}

	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView(), m.helpView())
}

func (m *ConfigView) sync() {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	helpHeight := lipgloss.Height(m.helpView())
	verticalMarginHeight := headerHeight + footerHeight + helpHeight

	if !m.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		offset := m.viewport.YOffset

		m.viewport = viewport.New(m.width, m.height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.SetContent(m.view)
		m.viewport.YOffset = offset
		m.ready = true
	} else {
		m.viewport.Width = m.width
		m.viewport.Height = m.height - verticalMarginHeight
	}
}

func (m *ConfigView) headerView() string {
	title := styleTitle().Render("Config: " + m.bus.Config().Path())
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *ConfigView) footerView() string {
	info := styleInfo().Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)) //nolint:mnd // not a magic number
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *ConfigView) helpView() string {
	return m.help.View(m.keys)
}
