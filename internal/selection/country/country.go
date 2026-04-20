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
	list       list.Model
	source     source
	fetching   bool
	totalCount int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		s.LightDark = lipgloss.LightDark(msg.IsDark())
		m.list.Styles.Title = s.Styles.Title
		return m, nil

	case tea.WindowSizeMsg:
		width, height := s.Styles.App.GetFrameSize()
		m.list.SetSize(msg.Width-width, msg.Height-height)
		return m, nil

	case sel.PopulatedMsg:
		m.list.StopSpinner()
		m.fetching = false
		m.source.incrementChunk()
		return m, m.list.SetItems(append(m.list.Items(), msg...))
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.list = newListModel
	if m.list.Paginator.OnLastPage() &&
		!m.fetching &&
		!m.IsFiltering() &&
		len(m.list.Items()) < m.totalCount {
		m.fetching = true
		cmds = append(cmds, m.Populate())
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	return tea.NewView(s.Styles.App.Render(m.list.View()))
}

func (m Model) Info() string {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		return i.ID
	}
	// TODO: should probably throw an error or something
	return ""
}

func (m *Model) Populate() tea.Cmd {
	m.fetching = true
	return tea.Batch(
		m.list.StartSpinner(),
		func() tea.Msg {
			// m.list.SetItems(append(m.list.Items(), m.Source.items()...))
			return sel.PopulatedMsg(m.source.items())
		},
	)
}

func (m *Model) Select(msg sel.Msg) {
	m.source.currentChunk = 0
	var items []list.Item
	m.list.SetItems(items)
	switch msg := msg.(type) {
	case sel.CountrySelectedMsg:
		code := msg.Code
		size := msg.Count
		m.totalCount = size
		if code == "ALL" {
			m.source = allStationSource()
		} else {
			m.source = stationsByCountryCodeSource(code)
		}
	}
}

func (m *Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m *Model) ViewLayer() *lipgloss.Layer {
	return lipgloss.NewLayer(s.Styles.App.Render(m.list.View()))
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
	m.fetching = false
	m.list.SetSpinner(spinner.Line)

	return m
}
