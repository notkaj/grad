package world

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = update
	d.ShortHelpFunc = listKeys.ShortHelp
	d.FullHelpFunc = listKeys.FullHelp
	return d
}

func update(msg tea.Msg, m *list.Model) tea.Cmd {
	// i, ok := m.SelectedItem().(item)
	// if !ok {
	// 	return nil
	// }

	switch b := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(b, listKeys.toggleSpinner):
			return m.ToggleSpinner()

		case key.Matches(b, listKeys.toggleTitleBar):
			v := !m.ShowTitle()
			m.SetShowTitle(v)
			m.SetShowFilter(v)
			m.SetFilteringEnabled(v)

		case key.Matches(b, listKeys.toggleStatusBar):
			m.SetShowStatusBar(!m.ShowStatusBar())

		case key.Matches(b, listKeys.togglePagination):
			m.SetShowPagination(!m.ShowPagination())

		case key.Matches(b, listKeys.toggleHelpMenu):
			m.SetShowHelp(!m.ShowHelp())

		}
	}
	return nil
}

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

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.toggleSpinner,
		d.toggleTitleBar,
		d.toggleStatusBar,
		d.togglePagination,
		d.toggleHelpMenu,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.toggleSpinner,
			d.toggleTitleBar,
			d.toggleStatusBar,
			d.togglePagination,
			d.toggleHelpMenu,
		},
	}
}
