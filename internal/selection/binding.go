package selection

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Choose    key.Binding
	Quickview key.Binding
	Back      key.Binding
}

var Keys = KeyMap{
	Choose: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
	Quickview: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "quick view"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "back"),
	),
}
