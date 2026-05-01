// Package selection provides screens for selecting radio stations
package selection

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/notkaj/grad/internal/card"
	"github.com/notkaj/grad/internal/selection/countries"
	"github.com/notkaj/grad/internal/selection/search"
	sel "github.com/notkaj/grad/internal/selection/selector"
	"github.com/notkaj/grad/internal/selection/stations"
)

type Model struct {
	width, height int
	countries     *countries.Model
	stations      *stations.Model
	search        *search.Model
	screen        sel.Selector
	playbackCard  card.PlaybackModel
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.screen.Init(), m.playbackCard.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.stations.Update(msg)
		m.countries.Update(msg)
		m.search.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.countries.Update(msg)
		m.stations.Update(msg)
		m.search.Update(msg)
		return m, nil
	case sel.CountrySelectedMsg:
		if msg.Code == "ALL" {
			m.screen = m.search
		} else {
			m.screen = m.stations
		}

	case tea.KeyPressMsg:
		if m.screen.IsFiltering() {
			_, cmd := m.screen.Update(msg)
			return m, cmd
		}
		switch {
		case key.Matches(msg, Keys.Back):
			switch m.screen.(type) {
			case *stations.Model, *search.Model:
				m.screen = m.countries
				return m, nil
			}
		}
	}
	_, screenCmd := m.screen.Update(msg)
	res, cardCmd := m.playbackCard.Update(msg)
	m.playbackCard = res.(card.PlaybackModel)
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
	countriesModel := countries.InitialModel()
	stationsModel := stations.InitialModel()
	searchModel := search.InitialModel()
	return &Model{
		countries:    &countriesModel,
		stations:     &stationsModel,
		search:       &searchModel,
		screen:       &countriesModel,
		playbackCard: card.InitialPlaybackModel(),
	}
}
