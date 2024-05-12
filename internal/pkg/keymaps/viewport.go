package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type ViewportKeymap struct {
	viewport.KeyMap
	Left key.Binding
	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (v ViewportKeymap) ShortHelp() []key.Binding {
	return []key.Binding{v.Help, v.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (v ViewportKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{v.Up, v.Down},                 // first column
		{v.PageUp, v.PageDown},         // second column
		{v.HalfPageUp, v.HalfPageDown}, // third column
		{v.Help, v.Quit},               // fourth column
		{v.Left},                       // fourth column
	}
}

func Viewport() ViewportKeymap {
	v := ViewportKeymap{
		KeyMap: viewport.KeyMap{
			PageDown: key.NewBinding(
				key.WithKeys("pgdown", " ", "f"),
				key.WithHelp("f / pg dn", ": page down"),
			),
			PageUp: key.NewBinding(
				key.WithKeys("pgup", "b"),
				key.WithHelp("b / pg up", ": page up"),
			),
			HalfPageUp: key.NewBinding(
				key.WithKeys("u", "ctrl+u"),
				key.WithHelp("u", ": ½ page up"),
			),
			HalfPageDown: key.NewBinding(
				key.WithKeys("d", "ctrl+d"),
				key.WithHelp("d", ": ½ page down"),
			),
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑ / k", ": up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓ / j", ": down"),
			),
		},
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", ": move left"),
		),
		Help: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", ": toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp("q / esc", ": quit"),
		),
	}

	return v
}
