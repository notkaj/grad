// Package chooser provides list of countries
package chooser

import (
	"sync"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styleMap struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

var (
	lightDark = lipgloss.LightDark(true)
	styles    = styleMap{
		app: lipgloss.
			NewStyle().
			Padding(1, 2).
			Margin(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#c4068b")),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#c4068b")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#c4068b"), lipgloss.Color("#c4068b"))),
	}
)

type item struct {
	title       string
	description string
	id          string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type keyMap struct {
	choose    key.Binding
	quickview key.Binding
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
}

type model struct {
	width, height int
	once          *sync.Once
	list          list.Model
	source        source
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

	i, ok := m.list.SelectedItem().(item)
	if ok {
		title := i.Title()
		id := i.id

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
				statusCmd := m.list.NewStatusMessage(styles.statusMessage.Render("You chose " + title))
				m.list.Title = string(m.source.category())
				return m, tea.Batch(itemsCmd, statusCmd)
			}

			if key.Matches(msg, keys.quickview) {
				return m, m.list.NewStatusMessage(styles.statusMessage.Render("quick view for " + id))
			}
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

func InitialModel() model {
	// Initialize the model and list.
	m := model{}
	lightDark = lipgloss.LightDark(true)
	m.source = defaultCountrySource()

	items := m.source.items()

	// Setup list.
	delegate := newItemDelegate()
	countryList := list.New(items, delegate, 0, 0)
	countryList.Title = string(m.source.category())
	countryList.Styles.Title = styles.title
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
