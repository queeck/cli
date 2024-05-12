package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
)

// DefaultKeymap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type DefaultKeymap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k DefaultKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k DefaultKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},     // first column
		{k.Left, k.Right},  // second column
		{k.Select, k.Quit}, // third column
		{k.Help},           // fourth column
	}
}

func Default() DefaultKeymap {
	return DefaultKeymap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", ": move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", ": move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", ": move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", ": move right"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter / space", ": select"),
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
}
