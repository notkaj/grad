// Package selection provides screens for selecting radio stations
package selection

import (
	tea "charm.land/bubbletea/v2"
	c "github.com/notkaj/grad/internal/selection/country"
	w "github.com/notkaj/grad/internal/selection/world"
)

type Root struct {
	world   *w.Model
	country *c.Model
	screen  tea.Model
}

func InitialModel() Root {
	worldModel := w.InitialModel()
	countryModel := c.InitialModel()
	return Root{
		world:   &worldModel,
		country: &countryModel,
		screen:  &worldModel,
	}
}
