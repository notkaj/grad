package world

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	sel "github.com/notkaj/grad/internal/selection/selector"
	g "gitlab.com/AgentNemo/goradios"
)

type source struct {
	order      g.Order
	reversed   bool
	filter     string
	includeAcc bool
}

func (c source) items() []list.Item {
	countries := g.FetchCountries()
	len := len(countries)
	items := make([]list.Item, len+1)
	sum := 0
	for i, country := range countries {
		desc := fmt.Sprintf("%d Stations", country.StationCount)
		items[i+1] = sel.NewItem(country.Name, desc, country.Code)
		sum += country.StationCount
	}

	allItem := sel.AllItem(sum)
	items[0] = allItem
	if c.includeAcc {
		return items
	}
	return items[1:]
}

func defaultCountrySource() source {
	return source{
		order:      g.OrderName,
		reversed:   false,
		includeAcc: true,
	}
}
