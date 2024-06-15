package keymaps

import "github.com/charmbracelet/bubbles/key"

type InputKeymap struct {
	Left   key.Binding
	Tab    key.Binding
	Quit   key.Binding
	Select key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k InputKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Left, k.Tab, k.Quit, k.Select,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k InputKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}

func Input() InputKeymap {
	return InputKeymap{
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", ": move left"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", ": autocomplete"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", ": select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc / ctrl+c", ": quit"),
		),
	}
}
