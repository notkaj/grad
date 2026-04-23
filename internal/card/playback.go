// Package card provides card layers for UI
package card

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/notkaj/grad/internal/animation"
	"github.com/notkaj/grad/internal/playback"
	sel "github.com/notkaj/grad/internal/selection/selector"
	s "github.com/notkaj/grad/internal/style"
)

type PlaybackModel struct {
	StationName    string
	Codec          string
	Status         string
	playingSpinner spinner.Model
	loadingSpinner spinner.Model
	Player         *playback.Player
}

func InitialPlaybackModel() PlaybackModel {
	s := spinner.New()
	s.Spinner = animation.Playing()
	return PlaybackModel{Status: "Nothing Selected", Player: playback.NewPlayer(), playingSpinner: s}
}

func (m PlaybackModel) Init() tea.Cmd {
	return m.playingSpinner.Tick
}

func (m PlaybackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playback.PlayerStartedMsg:
		m.Status = string(msg)
		return m, nil
	case playback.PlayerErrorMsg:
		m.Status = "Error: " + string(msg)
		return m, nil

	case sel.StationSelectedMsg:
		m.StationName = msg.Title
		m.Codec = msg.Codec
		url := msg.URL
		return m, m.Player.Play(url)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.togglePlayback):
			if m.Player == nil {
				return m, nil
			}
			if m.Player.IsPlaying() {
				m.Player.Stop()
			} else {
				return m, m.Player.Play(m.Player.URL)
			}
		case key.Matches(msg, keys.volumeUp):
			m.Player.SetVolume(m.Player.Volume() + 0.1)
			return m, nil
		case key.Matches(msg, keys.volumeDown):
			m.Player.SetVolume(m.Player.Volume() - 0.1)
			return m, nil
		case key.Matches(msg, keys.mute):
			m.Player.ToggleMute()
			return m, nil
		}

	}
	var playingCmd tea.Cmd
	var loadingCmd tea.Cmd
	m.playingSpinner, playingCmd = m.playingSpinner.Update(msg)
	m.loadingSpinner, loadingCmd = m.loadingSpinner.Update(msg)
	return m, tea.Batch(playingCmd, loadingCmd)
}

func (m PlaybackModel) View() tea.View {
	if m.StationName == "" {
		return tea.NewView("Select Station")
	}
	return tea.NewView("use Layer method")
}

func (m PlaybackModel) ViewLayer(parentWidth, parentHeight int) *lipgloss.Layer {
	height := s.Styles.Card.GetHeight()
	if m.StationName == "" {
		return lipgloss.NewLayer(
			s.Styles.Card.Render("Select Station!"),
		).
			X(parentWidth/2 - 10).
			Y(parentHeight/2 - height/2)
	}
	// what info do i want?
	// station name
	// home url
	// play status
	// - nothing selected
	// - loading
	// - playing
	// - stopped
	// - error
	// volume/mute
	// format/codec
	// location
	//
	header := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Align(lipgloss.Center).Render(m.StationName)
	spinner := m.playingSpinner.View()
	content := lipgloss.JoinVertical(lipgloss.Top, header, spinner)

	return lipgloss.NewLayer(
		s.Styles.Card.Render(content),
	).X(parentWidth/2 - 10).Y(parentHeight/2 - height/2)
}

// func (m PlaybackModel) volumeToPercent() float64 {
// 	// Using base 2, so 2^Volume * 100
// 	return math.Pow(2, m.Player.Volume()) * 100
// }
