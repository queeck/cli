package view

import "github.com/charmbracelet/lipgloss"

func styleTitle() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
}

func styleInfo() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	b.Left = "┤"
	return styleTitle().Copy().BorderStyle(b)
}
