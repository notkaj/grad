// Package selection provides screens for selecting radio stations
package selection

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/notkaj/grad/internal/card"
	c "github.com/notkaj/grad/internal/selection/country"
	sel "github.com/notkaj/grad/internal/selection/selector"
	w "github.com/notkaj/grad/internal/selection/world"
)

type Model struct {
	width, height int
	world         *w.Model
	country       *c.Model
	screen        sel.Selector
	playbackCard  card.PlaybackModel
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.screen.Init(), m.playbackCard.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.country.Update(msg)
		m.world.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.world.Update(msg)
		m.country.Update(msg)
		return m, nil
	case sel.CountrySelectedMsg:
		m.screen = m.country

	case tea.KeyPressMsg:
		if m.screen.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, Keys.Back):
			switch m.screen.(type) {
			case *c.Model:
				m.screen = m.world
				return m, nil
			}
		}
	}
	_, screenCmd := m.screen.Update(msg)
	c, cardCmd := m.playbackCard.Update(msg)
	m.playbackCard = c.(card.PlaybackModel)
	return m, tea.Batch(screenCmd, cardCmd)
}

func (m *Model) View() tea.View {
	backdrop := m.screen.ViewLayer()
	playback := m.playbackCard.
		ViewLayer(m.width, m.height)

	comp := lipgloss.NewCompositor(
		backdrop,
		playback,
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
		world:        &worldModel,
		country:      &countryModel,
		screen:       &worldModel,
		playbackCard: card.InitialPlaybackModel(),
	}
}
