// Package selector provides interface Selector
package selector

import tea "charm.land/bubbletea/v2"

type Selector interface {
	tea.Model
	ID() string
	Populate() tea.Cmd
	Select(Msg)
	IsFiltering() bool
}

type PopulatedMsg string

type Msg any

type (
	CountryCodeMsg string
)
