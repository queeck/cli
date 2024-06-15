package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/config"
	"github.com/queeck/cli/internal/services/templates"
	"github.com/queeck/cli/internal/services/templates/texts"
)

const (
	Code = "view"
)

type Templates interface {
	templates.RendererCommon
	templates.RendererConfigView
}

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	keymap    keymaps.ViewportKeymap
	help      help.Model
	viewport  viewport.Model
	templates Templates
	parent    func(commands.Command) commands.Command
	config    config.Config
	view      string
	quitting  bool
	ready     bool
	width     int
	height    int
}

func New(bus commands.Bus) commands.Command {
	return &Model{
		keymap:    keymaps.Viewport(),
		help:      help.New(),
		view:      bus.Config().View(),
		templates: bus.Templates(),
		parent:    bus.Parent,
		config:    bus.Config(),
		ready:     false,
	}
}

func (m *Model) Code() string {
	return Code
}

func (m *Model) Commands() []commands.Variant {
	return nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = m.width
		m.sync()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.parent(m), nil
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.ready = false
			m.sync()
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, tea.Batch(cmd)
}

func (m *Model) View() string {
	if m.quitting {
		return m.templates.RenderCommonQuit()
	}

	if !m.ready {
		return texts.ConfigViewInitializing
	}

	return m.templates.RenderConfigViewScreen(m.headerView(), m.viewport.View(), m.footerView(), m.helpView())
}

func (m *Model) sync() {
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

func (m *Model) headerView() string {
	configPath := m.config.Path()
	headerTitle := m.templates.RenderConfigViewHeaderTitle(configPath)
	title := styleTitle().Render(headerTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	scrollPercent := m.viewport.ScrollPercent()
	footerInfo := m.templates.RenderConfigViewFooterInfo(scrollPercent)
	info := styleInfo().Render(footerInfo)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *Model) helpView() string {
	return m.help.View(m.keymap)
}
