// Package world provides list of countries
package world

import (
	"sync"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	c "github.com/notkaj/grad/internal/selection/common"
	s "github.com/notkaj/grad/internal/style"
)

type keyMap struct {
	choose    key.Binding
	quickview key.Binding
	back      key.Binding
}

var keys = keyMap{
	choose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
	quickview: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "quick view"),
	),
	back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "back"),
	),
}

type Model struct {
	width, height int
	once          *sync.Once
	list          list.Model
	source        source
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *Model) updateListProperties() {
	// Update list size.
	h, v := s.Styles.App.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	// Update the model and list styles.
	m.list.Styles.Title = s.Styles.Title
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		s.LightDark = lipgloss.LightDark(msg.IsDark())
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	i, ok := m.list.SelectedItem().(c.Item)
	if ok {
		title := i.Title()
		id := i.ID

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			// Don't match any of the keys below if we're actively filtering.
			if m.list.FilterState() == list.Filtering {
				break
			}

			if key.Matches(msg, keys.choose) {
				if m.source.category() == stations {
					return m, nil
				}
				m.source = defaultStationSource(id)
				itemsCmd := m.list.SetItems(m.source.items())
				statusCmd := m.list.NewStatusMessage(s.Styles.StatusMessage.Render("You chose " + title))
				m.list.Title = string(m.source.category())
				return m, tea.Batch(itemsCmd, statusCmd)
			}

			if key.Matches(msg, keys.back) {
				if m.source.category() == countries {
					return m, nil
				}
				m.source = defaultCountrySource()
				itemsCmd := m.list.SetItems(m.source.items())
				m.list.Title = string(m.source.category())
				return m, itemsCmd
			}

			if key.Matches(msg, keys.quickview) {
				return m, m.list.NewStatusMessage(s.Styles.StatusMessage.Render("quick view for " + id))
			}
		}
	}
	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	v := tea.NewView(s.Styles.App.Render(m.list.View()))
	v.AltScreen = true
	return v
}

func InitialModel() Model {
	// Initialize the model and list.
	m := Model{}
	s.LightDark = lipgloss.LightDark(true)
	m.source = defaultCountrySource()

	items := m.source.items()

	// Setup list.
	delegate := newItemDelegate()
	countryList := list.New(items, delegate, 0, 0)
	countryList.Title = string(m.source.category())
	countryList.Styles.Title = s.Styles.Title
	countryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.choose,
			keys.quickview,
		}
	}

	m.list = countryList
	m.list.SetSpinner(spinner.Dot)

	return m
}
