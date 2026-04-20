// Package selection provides screens for selecting radio stations
package selection

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/notkaj/grad/internal/card"
	"github.com/notkaj/grad/internal/playback"
	c "github.com/notkaj/grad/internal/selection/country"
	sel "github.com/notkaj/grad/internal/selection/selector"
	w "github.com/notkaj/grad/internal/selection/world"
)

type Model struct {
	width, height int
	world         *w.Model
	country       *c.Model
	screen        sel.Selector
	player        *playback.Player
	playbackCard  card.PlaybackModel
}

func (m *Model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playback.PlayerStartedMsg:
		m.playbackCard.Status = string(msg)
		return m, nil
	case playback.PlayerErrorMsg:
		m.playbackCard.Status = "Error: " + string(msg)
		return m, nil
	case tea.BackgroundColorMsg:
		m.country.Update(msg)
		m.world.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.world.Update(msg)
		m.country.Update(msg)
		return m, nil
	case tea.KeyPressMsg:
		if m.screen.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, Keys.Select):
			switch t := m.screen.(type) {
			case *w.Model:
				id, name, size := t.SelectionInfo()
				m.country.Select(sel.CountrySelectedMsg{Code: id, Name: name, Count: size})
				m.screen = m.country
				return m, m.screen.Populate()
			case *c.Model:
				name, url, codec := t.SelectionInfo()
				m.playbackCard.StationName = name
				m.playbackCard.Codec = codec
				m.playbackCard.Status = "Loading..."
				return m, m.player.Play(url)
			}

		case key.Matches(msg, Keys.Back):
			switch m.screen.(type) {
			case *c.Model:
				m.screen = m.world
				return m, nil
			}
		case key.Matches(msg, Keys.TogglePlayback):
			if m.player == nil {
				return m, nil
			}
			if m.player.IsPlaying() {
				m.player.Stop()
			} else {
				return m, m.player.Play(m.player.URL)
			}
		}
	}
	_, cmd := m.screen.Update(msg)
	return m, cmd
}

func (m *Model) View() tea.View {
	backdrop := m.screen.ViewLayer()
	playbackLayer := m.playbackCard.Layer(m.width, m.height)

	comp := lipgloss.NewCompositor(
		backdrop,
		playbackLayer,
	)
	var view tea.View
	view.SetContent(comp.Render())
	view.AltScreen = true
	return view
}

func InitialModel() *Model {
	worldModel := w.InitialModel()
	countryModel := c.InitialModel()
	return &Model{
		world:   &worldModel,
		country: &countryModel,
		screen:  &worldModel,
		player:  playback.NewPlayer(),
	}
}
