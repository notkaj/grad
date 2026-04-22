package selection

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Select         key.Binding
	Quickview      key.Binding
	Back           key.Binding
	TogglePlayback key.Binding
	VolumeUp       key.Binding
	VolumeDown     key.Binding
	Mute           key.Binding
}

var Keys = KeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
	Quickview: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "quick view"),
	),
	Back: key.NewBinding(
		key.WithKeys("backspace", "ctrl+h", "esc"),
		key.WithHelp("backspace", "back"),
	),
	TogglePlayback: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "toggle playback"),
	),
	VolumeUp: key.NewBinding(
		key.WithKeys("+", "="),
		key.WithHelp("+", "volume up"),
	),
	VolumeDown: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "volume down"),
	),
	Mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "mute"),
	),
}
