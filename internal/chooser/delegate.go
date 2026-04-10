package chooser

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = update
	d.ShortHelpFunc = keys.ShortHelp
	d.FullHelpFunc = keys.FullHelp
	return d
}

func update(msg tea.Msg, m *list.Model) tea.Cmd {
	i, ok := m.SelectedItem().(item)
	if !ok {
		return nil
	}

	title := i.Title()

	switch b := msg.(type) {
	case tea.KeyPressMsg:

		if key.Matches(b, keys.choose) {
			return m.NewStatusMessage(styles.statusMessage.Render("You chose " + title))
		}

		if key.Matches(b, keys.quickview) {
			return m.NewStatusMessage(styles.statusMessage.Render("quick view: "))
		}
	}
	return nil
}

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

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.quickview,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.quickview,
		},
	}
}
