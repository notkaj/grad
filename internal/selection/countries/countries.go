// Package countries provides list of countries
package countries

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	sel "github.com/notkaj/grad/internal/selection/selector"
	s "github.com/notkaj/grad/internal/style"
)

type Model struct {
	list          list.Model
	source        source
	width, height int
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
		fw, fh := s.Styles.App.GetFrameSize()
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(m.width-fw, m.height-fh)
		return m, nil

	case sel.PopulatedMsg:
		m.list.SetItems(msg)
		m.list.StopSpinner()

	case tea.KeyPressMsg:
		if m.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, sel.Keys.Select):
			return m, func() tea.Msg {
				code, name := m.selectionInfo()
				return sel.CountrySelectedMsg{Code: code, Name: name}
			}
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() tea.View {
	return tea.NewView(
		s.Styles.App.
			Render(m.list.View()),
	)
}

func (m *Model) selectionInfo() (string, string) {
	if i, ok := m.list.SelectedItem().(item); ok {
		return i.id, i.title
	}
	return "", ""
}

func (m *Model) Populate() tea.Cmd {
	return tea.Batch(
		m.list.StartSpinner(),
		func() tea.Msg {
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
	return lipgloss.
		NewLayer(
			s.Styles.App.
				Render(m.list.View()),
		)
}

func InitialModel() Model {
	// Initialize the model and list.
	m := Model{}
	s.LightDark = lipgloss.LightDark(true)
	m.source = defaultCountrySource()

	// Setup list.
	delegate := newItemDelegate()
	delegate.Styles.SelectedTitle = s.Styles.SeletectedTitle
	delegate.Styles.SelectedDesc = s.Styles.SelectedDesc
	delegate.Styles.NormalTitle = s.Styles.ItemTitle
	delegate.Styles.NormalDesc = s.Styles.ItemDesc
	countryList := list.New(nil, delegate, 0, 0)
	countryList.Title = "Countries"
	countryList.Styles.Title = s.Styles.Title

	m.list = countryList
	m.list.SetSpinner(spinner.Line)

	return m
}
