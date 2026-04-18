// Package selector provides interface Selector
package selector

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type Selector interface {
	tea.Model
	Populate() tea.Cmd
	Select(Msg)
	IsFiltering() bool
}

type PopulatedMsg []list.Item

type Msg any

type CountrySelectedMsg struct {
	Code  string
	Count int
}
