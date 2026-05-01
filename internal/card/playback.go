// Package card provides card layers for UI
package card

import (
	"fmt"
	"math"

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
	StationName     string
	Codec           string
	Status          string
	Homepage        string
	playbackSpinner spinner.Model
	loadingSpinner  spinner.Model
	Player          *playback.Player
	state           state
}

func InitialPlaybackModel() PlaybackModel {
	playbackSpinner := spinner.New()
	playbackSpinner.Spinner = animation.Playing()
	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Dot
	return PlaybackModel{
		Status:          "Nothing Selected",
		Player:          playback.NewPlayer(),
		playbackSpinner: playbackSpinner,
		loadingSpinner:  loadingSpinner,
		state:           noSelection,
	}
}

func (m PlaybackModel) Init() tea.Cmd {
	return tea.Batch(m.playbackSpinner.Tick, m.loadingSpinner.Tick)
}

func (m PlaybackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playback.PlayerStartedMsg:
		m.Status = string(msg)
		m.state = playing
		return m, nil
	case playback.PlayerErrorMsg:
		m.Status = "Error: " + string(msg)
		m.state = error
		return m, nil

	case sel.StationSelectedMsg:
		m.StationName = msg.Title
		m.Codec = msg.Codec
		m.Homepage = msg.Homepage
		m.state = loading
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
				m.state = stopped
			} else {
				m.state = playing
				return m, m.Player.Play(m.Player.URL)
			}
		case key.Matches(msg, keys.volumeUp):
			m.Player.SetVolume(m.Player.Volume() + 0.1)
			return m, nil
		case key.Matches(msg, keys.volumeDown):
			m.Player.SetVolume(m.Player.Volume() - 0.1)
			return m, nil
		case key.Matches(msg, keys.mute):
			m.state = muted
			m.Player.ToggleMute()
			return m, nil
		}

	}
	var playingCmd tea.Cmd
	var loadingCmd tea.Cmd
	m.playbackSpinner, playingCmd = m.playbackSpinner.Update(msg)
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
	fullWidth := s.Styles.Card.GetWidth()
	fullHeight := s.Styles.Card.GetHeight()
	frameWidth, frameHeight := s.Styles.Card.GetFrameSize()
	width, height := fullWidth-frameWidth, fullHeight-frameHeight
	if m.state == noSelection {
		// return lipgloss.NewLayer(
		// 	s.Styles.Card.Render("Select Station!"),
		// ).
		// 	X(parentWidth/2 - 10).
		// 	Y(parentHeight/2 - height/2)
		return lipgloss.NewLayer("")
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
	thirdWidth := width / 3
	volumePercent := fmt.Sprintf("Volume: %.0f%%", m.volumeToPercent())
	header := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center).
		Render(m.StationName)
	footer := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().
			Width(thirdWidth).
			Align(lipgloss.Left).
			Render(m.Codec),
		lipgloss.NewStyle().
			Width(thirdWidth).
			Align(lipgloss.Center).
			Render(volumePercent),
		lipgloss.NewStyle().
			Width(width-2*thirdWidth).
			Align(lipgloss.Right).
			Render(m.Status),
	)

	playbackVisual := "██"

	if m.state == playing {
		playbackVisual = m.playbackSpinner.View()
	}
	if m.state == loading {
		playbackVisual = m.loadingSpinner.View()
	}

	main := lipgloss.NewStyle().
		Width(width).
		Height(height - lipgloss.Height(header) - lipgloss.Height(footer)).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(playbackVisual + "\n" + m.Homepage)

	content := lipgloss.
		JoinVertical(
			lipgloss.Top,
			header,
			main,
			footer,
		)

	return lipgloss.NewLayer(
		s.Styles.Card.Render(content),
	).X(parentWidth/2 - 10).Y(parentHeight/2 - height/2)
}

func (m PlaybackModel) volumeToPercent() float64 {
	// Using base 2, so 2^Volume * 100
	return math.Pow(2, m.Player.Volume()) * 100
}

type state string

const (
	playing     = state("playing")
	loading     = state("loading")
	error       = state("error")
	noSelection = state("nothing selected")
	muted       = state("muted")
	stopped     = state("stopped")
)
