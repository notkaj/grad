package card

import "charm.land/bubbles/v2/key"

type keyMap struct {
	togglePlayback key.Binding
	volumeUp       key.Binding
	volumeDown     key.Binding
	mute           key.Binding
}

var keys = keyMap{
	togglePlayback: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "toggle playback"),
	),
	volumeUp: key.NewBinding(
		key.WithKeys("+", "="),
		key.WithHelp("+", "volume up"),
	),
	volumeDown: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "volume down"),
	),
	mute: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "mute"),
	),
}
