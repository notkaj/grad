// Package world provides list of countries
package world

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	sel "github.com/notkaj/grad/internal/selection/selector"
	s "github.com/notkaj/grad/internal/style"
)

type Model struct {
	list   list.Model
	source source
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.Populate(),
		tea.RequestBackgroundColor,
	)
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
		m.list.SetItems(msg)
		m.list.StopSpinner()
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	return tea.NewView(s.Styles.App.Render(m.list.View()))
}

func (m *Model) SelectionInfo() (string, int) {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		return i.ID, i.Count
	}
	return "", 0
}

func (m *Model) Populate() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		func() tea.Msg {
			// m.list.SetItems(m.source.items())
			return sel.PopulatedMsg(m.source.items())
		},
	)
}

func (m *Model) Select(msg sel.Msg) {
}

func (m *Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m *Model) ViewLayer() *lipgloss.Layer {
	return lipgloss.NewLayer(s.Styles.App.Render(m.list.View()))
}

func InitialModel() Model {
	// Initialize the model and list.
	m := Model{}
	s.LightDark = lipgloss.LightDark(true)
	m.source = defaultCountrySource()

	// Setup list.
	delegate := newItemDelegate()
	countryList := list.New(nil, delegate, 0, 0)
	countryList.Title = "Countries"
	countryList.Styles.Title = s.Styles.Title

	m.list = countryList
	m.list.SetSpinner(spinner.Line)

	return m
}
