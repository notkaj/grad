// Package style provides lipgloss styles for ui elements
package style

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

type styleMap struct {
	App             lipgloss.Style
	Title           lipgloss.Style
	StatusMessage   lipgloss.Style
	ItemTitle       lipgloss.Style
	ItemDesc        lipgloss.Style
	SeletectedTitle lipgloss.Style
	SelectedDesc    lipgloss.Style
	Card            lipgloss.Style
}

var (
	foreground = charmtone.Salt
	base       = charmtone.Coral
	selected   = charmtone.Tang
	LightDark  = lipgloss.LightDark(true)
	Styles     = styleMap{
		App: lipgloss.
			NewStyle().
			Padding(1, 2).
			Margin(1).
			// BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(base),
		Title: lipgloss.NewStyle().
			Foreground(charmtone.Charcoal).
			Background(base).
			Padding(0, 1),
		StatusMessage: lipgloss.NewStyle().
			Foreground(LightDark(base, base)),
		ItemTitle:       lipgloss.NewStyle().Foreground(foreground),
		ItemDesc:        lipgloss.NewStyle().Foreground(charmtone.Squid),
		SeletectedTitle: lipgloss.NewStyle().Foreground(selected),
		SelectedDesc:    lipgloss.NewStyle().Foreground(selected),
		Card: lipgloss.NewStyle().
			Width(40).
			Height(20).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Lilac).
			BorderForegroundBlend(
				charmtone.Cherry,
				charmtone.Charple,
				charmtone.Guac,
				charmtone.Charple,
				charmtone.Sriracha,
			).
			Foreground(foreground).
			Align(lipgloss.Center, lipgloss.Center),
	}
)
