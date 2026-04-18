// Package selection provides screens for selecting radio stations
package selection

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	c "github.com/notkaj/grad/internal/selection/country"
	sel "github.com/notkaj/grad/internal/selection/selector"
	w "github.com/notkaj/grad/internal/selection/world"
)

type Model struct {
	world   *w.Model
	country *c.Model
	screen  sel.Selector
}

func (m *Model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		// s.LightDark = lipgloss.LightDark(msg.IsDark())
		// m.updateListProperties()
		// return m, nil
		m.country.Update(msg)
		m.world.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		// m.width, m.height = msg.Width, msg.Height
		// m.updateListProperties()
		// return m, nil
		m.world.Update(msg)
		m.country.Update(msg)
	case tea.KeyPressMsg:
		if m.screen.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, Keys.Select):
			switch t := m.screen.(type) {
			case *w.Model:
				id := t.ID()
				m.country.Select(sel.CountryCodeMsg(id))
				m.screen = m.country
				return m, m.screen.Populate()
			}
		case key.Matches(msg, Keys.Back):
			switch m.screen.(type) {
			case *c.Model:
				m.screen = m.world
				return m, nil
			}
		}
	}
	_, cmd := m.screen.Update(msg)
	return m, cmd
}

func (m *Model) View() tea.View {
	v := m.screen.View()
	v.AltScreen = true
	return v
}

func InitialModel() *Model {
	worldModel := w.InitialModel()
	countryModel := c.InitialModel()
	return &Model{
		world:   &worldModel,
		country: &countryModel,
		screen:  &worldModel,
	}
}
