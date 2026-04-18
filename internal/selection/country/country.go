// Package country provides selector for stations
package country

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	sel "github.com/notkaj/grad/internal/selection/selector"
)

type Model struct {
	list   list.Model
	Source source
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case sel.PopulatedMsg:
		m.list.StopSpinner()
		return m, nil
	}
	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView("not yet implemented")
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
