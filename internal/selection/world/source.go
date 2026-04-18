package world

import (
	"fmt"

	"charm.land/bubbles/v2/list"
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
	sum := int(0)
	for i, country := range countries {
		desc := fmt.Sprintf("%d Stations", country.StationCount)
		items[i+1] = NewItem(country.Name, desc, country.Code, int(country.StationCount))
		sum += country.StationCount
	}

	items[0] = AllItem(sum)
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
