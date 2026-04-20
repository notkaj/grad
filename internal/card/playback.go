// Package card provides cards/layers to display on top of the main UI
package card

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

type PlaybackModel struct {
	StationName string
	Status      string
}

func (m PlaybackModel) Layer(width, height int) *lipgloss.Layer {
	content := m.StationName
	if m.Status != "" {
		content += "\n\n" + m.Status
	}
	if content == "" {
		content = "No station selected"
	}

	cardWidth := 40
	cardHeight := 20

	return lipgloss.NewLayer(
		lipgloss.NewStyle().
			Width(cardWidth).
			Height(cardHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Damson).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content),
	).X(width / 2).Y(height/2 - cardHeight/2)
}
