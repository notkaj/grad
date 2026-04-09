package main

import (
	"fmt"
	"sync"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	g "gitlab.com/AgentNemo/goradios"
)

type styleMap struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

var (
	lightDark = lipgloss.LightDark(true)
	styles    = styleMap{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),
	}
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

var listKeys = listKeyMap{
	toggleSpinner: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle spinner"),
	),
	toggleTitleBar: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "toggle title"),
	),
	toggleStatusBar: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "toggle status"),
	),
	togglePagination: key.NewBinding(
		key.WithKeys("P"),
		key.WithHelp("P", "toggle pagination"),
	),
	toggleHelpMenu: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "toggle help"),
	),
}

type model struct {
	width, height int
	once          *sync.Once
	list          list.Model
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *model) updateListProperties() {
	// Update list size.
	h, v := styles.app.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.list.Styles.Title = styles.title
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		lightDark = lipgloss.LightDark(msg.IsDark())
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, listKeys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, listKeys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, listKeys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, listKeys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, listKeys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	v := tea.NewView(styles.app.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func initialModel() model {
	// Initialize the model and list.
	m := model{}
	lightDark = lipgloss.LightDark(true)

	// Make initial list of items.
	// var itemGenerator randomItemGenerator
	// const numItems = 24
	// items := make([]list.Item, numItems)
	// for i := range numItems {
	// 	items[i] = itemGenerator.next()
	// }

	countries := g.FetchCountries()
	items := make([]list.Item, len(countries))
	for i, c := range countries {
		items[i] = item{title: c.Name, description: fmt.Sprintf("Stations: %d", c.StationCount)}
	}

	// Setup list.
	delegate := newItemDelegate()
	countryList := list.New(items, delegate, 0, 0)
	countryList.Title = "Countries"
	countryList.Styles.Title = styles.title
	countryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	m.list = countryList

	return m
}
