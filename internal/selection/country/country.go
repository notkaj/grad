// Package country provides selector for stations
package country

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel

	return m, cmd
}

func (m *Model) View() tea.View {
	return tea.NewView(s.Styles.App.Render(m.list.View()))
}

func (m *Model) ID() string {
	item, ok := m.list.SelectedItem().(sel.Item)
	if ok {
		return item.ID
	}
	return ""
}

func (m *Model) Populate() tea.Cmd {
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

func (m *Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

func InitialModel() Model {
	m := Model{}
	s.LightDark = lipgloss.LightDark(true)

	// Setup list.
	delegate := newItemDelegate()
	stationsList := list.New(nil, delegate, 0, 0)
	stationsList.Title = "Stations"
	stationsList.Styles.Title = s.Styles.Title

	m.list = stationsList
	m.list.SetSpinner(spinner.Line)

	return m
}
