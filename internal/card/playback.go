// Package card provides card layers for UI
package card

import (
	"charm.land/lipgloss/v2"
	s "github.com/notkaj/grad/internal/style"
)

type PlaybackModel struct {
	StationName string
	Codec       string
	Status      string
}

func (m PlaybackModel) Layer(parentWidth, parentHeight int) *lipgloss.Layer {
	content := m.StationName
	if m.Codec != "" {
		content += "\n" + m.Codec
	}
	if m.Status != "" {
		content += "\n\n" + m.Status
	}
	if content == "" {
		content = "No station selected"
	}

	height := s.Styles.Card.GetHeight()

	return lipgloss.NewLayer(
		s.Styles.Card.Render(content),
	).X(parentWidth/2 - 5).Y(parentHeight/2 - height/2)
}
