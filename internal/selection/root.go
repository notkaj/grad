// Package selection provides screens for selecting radio stations
package selection

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
	c "github.com/notkaj/grad/internal/selection/country"
	sel "github.com/notkaj/grad/internal/selection/selector"
	w "github.com/notkaj/grad/internal/selection/world"
)

type Model struct {
	width, height int
	world         *w.Model
	country       *c.Model
	screen        sel.Selector
}

func (m *Model) Init() tea.Cmd {
	return m.screen.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.country.Update(msg)
		m.world.Update(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.world.Update(msg)
		m.country.Update(msg)
		return m, nil
	case tea.KeyPressMsg:
		if m.screen.IsFiltering() {
			break
		}
		switch {
		case key.Matches(msg, Keys.Select):
			switch t := m.screen.(type) {
			case *w.Model:
				id, size := t.SelectionInfo()
				m.country.Select(sel.CountrySelectedMsg{Code: id, Count: size})
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
	backdrop := m.screen.ViewLayer()
	card := lipgloss.NewLayer(
		lipgloss.NewStyle().
			Width(40).
			Height(20).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(charmtone.Damson).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Test"),
	).X(50).Y(20)

	comp := lipgloss.NewCompositor(
		backdrop,
		card,
	)
	var view tea.View
	view.SetContent(comp.Render())
	view.AltScreen = true
	return view
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
