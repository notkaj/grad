// Package style provides lipgloss styles for ui elements
package style

import "charm.land/lipgloss/v2"

type styleMap struct {
	App           lipgloss.Style
	Title         lipgloss.Style
	StatusMessage lipgloss.Style
}

var (
	LightDark = lipgloss.LightDark(true)
	Styles    = styleMap{
		App: lipgloss.
			NewStyle().
			Padding(1, 2).
			Margin(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#c4068b")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#c4068b")).
			Padding(0, 1),
		StatusMessage: lipgloss.NewStyle().
			Foreground(LightDark(lipgloss.Color("#c4068b"), lipgloss.Color("#c4068b"))),
	}
)
