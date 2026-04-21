// Package card provides card layers for UI
package card

import (
	"fmt"
	"math"

	"charm.land/lipgloss/v2"
	s "github.com/notkaj/grad/internal/style"
)

type PlaybackModel struct {
	StationName string
	Codec       string
	Status      string
	Volume      float64
	IsMuted     bool
}

func (m PlaybackModel) VolumeToPercent() float64 {
	// Using base 2, so 2^Volume * 100
	return math.Pow(2, m.Volume) * 100
}

func (m PlaybackModel) Layer(parentWidth, parentHeight int) *lipgloss.Layer {
	content := m.StationName
	if m.Codec != "" {
		content += "\n" + m.Codec
	}
	if m.Status != "" {
		content += "\n\n" + m.Status
	}

	vol := fmt.Sprintf("\nVolume: %.0f%%", m.VolumeToPercent())
	if m.IsMuted {
		vol += " (Muted)"
	}
	content += vol

	if content == "" {
		content = "No station selected"
	}

	height := s.Styles.Card.GetHeight()

	return lipgloss.NewLayer(
		s.Styles.Card.Render(content),
	).X(parentWidth/2 - 5).Y(parentHeight/2 - height/2)
}
