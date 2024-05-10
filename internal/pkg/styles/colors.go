package styles

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
)

func ColorForegroundSubtle(text string) string {
	return StyleForeground(text, subtle)
}

func ColorForegroundHighlight(text string) string {
	return StyleForeground(text, highlight)
}

func ColorForegroundSpecial(text string) string {
	return StyleForeground(text, special)
}

func StyleForeground(text string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().Foreground(color).Render(text)
}
