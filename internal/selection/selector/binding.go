package selector

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Select key.Binding
	Back   key.Binding
}

var Keys = keyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "ctrl+h", "esc"),
		key.WithHelp("backspace", "back"),
	),
}
