// Package country provides selector for stations
package country

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	sel "github.com/notkaj/grad/internal/selection/selector"
	s "github.com/notkaj/grad/internal/style"
)

type Model struct {
	width, height int
	list          list.Model
	Source        source
}

func (m *Model) updateListProperties() {
	// Update list size.
	h, v := s.Styles.App.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.list.Styles.Title = s.Styles.Title
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		s.LightDark = lipgloss.LightDark(msg.IsDark())
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil

	case sel.PopulatedMsg:
		m.list.StopSpinner()
		return m, nil
	}
	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView(s.Styles.App.Render(m.list.View()))
}

func (m Model) ID() string {
	item, ok := m.list.SelectedItem().(sel.Item)
	if ok {
		return item.ID
	}
	return ""
}

func (m Model) Populate() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		func() tea.Msg {
			m.list.SetItems(m.Source.items())
			return sel.PopulatedMsg("Stations Populated")
		},
	)
}

func (m *Model) Select(msg sel.Msg) {
	switch msg := msg.(type) {
	case sel.CountryCodeMsg:
		code := string(msg)
		if code == "ALL" {
			m.Source = allStationSource()
		} else {
			m.Source = stationsByCountryCodeSource(code)
		}
	}
}

func InitialModel() Model {
	return Model{
		list: list.Model{},
	}
}
