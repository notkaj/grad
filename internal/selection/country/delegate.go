package country

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
	switch b := msg.(type) {
	case tea.KeyPressMsg:
		switch {
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
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

var listKeys = listKeyMap{
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

func (d listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.toggleStatusBar,
		d.togglePagination,
		d.toggleHelpMenu,
	}
}

func (d listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.toggleStatusBar,
			d.togglePagination,
			d.toggleHelpMenu,
		},
	}
}
