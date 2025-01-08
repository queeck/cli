package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/queeck/cli/internal/pkg/keymaps"
	"github.com/queeck/cli/internal/services/commands"
	"github.com/queeck/cli/internal/services/commands/models"
)

const (
	Code = "view"
)

var _ commands.Command = &Model{} // check for interface compatibility

type Model struct {
	bus      commands.Bus
	keymap   keymaps.ViewportKeymap
	help     help.Model
	viewport viewport.Model
	quitting bool
	ready    bool
}

func New(bus commands.Bus) commands.Command {
	return &Model{
		bus:    bus,
		keymap: keymaps.Viewport(),
		help:   help.New(),
		ready:  false,
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

func (m *Model) OnUpdateWindowSize(msg tea.WindowSizeMsg) {
	m.bus.UpdateWindowSize(msg.Width, msg.Height)
	m.help.Width = msg.Width
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.OnUpdateWindowSize(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Left):
			return m.bus.Parent(m), nil
		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.ready = false
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, tea.Batch(cmd)
}

func (m *Model) View() string {
	m.sync()

	if m.quitting {
		return m.bus.Templates().Render(models.TemplateQuit)
	}

	if !m.ready {
		return m.bus.Templates().Render(templateInitializing)
	}

	return m.bus.Templates().Render(templateScreen,
		"header", m.headerView(),
		"viewport", m.viewport.View(),
		"footer", m.footerView(),
		"help", m.helpView(),
	)
}

func (m *Model) sync() {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	helpHeight := lipgloss.Height(m.helpView())
	verticalMarginHeight := headerHeight + footerHeight + helpHeight

	height := 0
	width := 0
	if window := m.bus.Window(); window != nil {
		height = window.Height()
		width = window.Width()
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			offset := m.viewport.YOffset

			m.viewport = viewport.New(width, height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.bus.Config().View())
			m.viewport.YOffset = offset
			m.ready = true
		}
	}

	m.viewport.Width = width
	m.viewport.Height = height - verticalMarginHeight
}

func (m *Model) headerView() string {
	configPath := m.bus.Config().Path()
	headerTitle := m.bus.Templates().Render(templateHeaderTitle, "path", configPath)
	title := styleTitle().Render(headerTitle)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	scrollPercent := m.viewport.ScrollPercent()
	footerInfo := m.bus.Templates().Render(templateFooterInfo,
		"scrollPercentage", fmt.Sprintf("%3.f", scrollPercent*100), //nolint:mnd // not a magic number
	)
	info := styleInfo().Render(footerInfo)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m *Model) helpView() string {
	return m.help.View(m.keymap)
}
