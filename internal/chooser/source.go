package chooser

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	g "gitlab.com/AgentNemo/goradios"
)

type source interface {
	items() []list.Item
	category() category
}

type category string

const (
	countries category = "Countries"
	stations  category = "Stations"
)

type countrySource struct {
	order      g.Order
	reversed   bool
	filter     string
	includeAcc bool
}

func (c countrySource) items() []list.Item {
	countries := g.FetchCountries()
	len := len(countries)
	items := make([]list.Item, len+1)
	sum := 0
	for i, country := range countries {
		items[i+1] = item{
			title:       country.Name,
			description: fmt.Sprintf("Station %d", country.StationCount),
			id:          country.Code,
		}
		sum += country.StationCount
	}
	items[0] = item{
		title:       "All",
		description: fmt.Sprintf("Stations %d", sum),
		id:          "ALL",
	}
	if c.includeAcc {
		return items
	}
	return items[1:]
}

func (countrySource) category() category {
	return countries
}

func defaultCountrySource() countrySource {
	return countrySource{
		order:      g.OrderName,
		reversed:   false,
		includeAcc: true,
	}
}

type stationSource struct {
	order    g.StationsOrder
	by       g.StationsBy
	term     string
	reversed bool
	filter   string
}

func (s stationSource) items() []list.Item {
	stations := g.FetchStationsDetailed(s.by, s.term, s.order, s.reversed, 0, 0, false)
	len := len(stations)
	items := make([]list.Item, len)
	for i, station := range stations {
		items[i] = item{
			title:       station.Name,
			description: station.URL,
			id:          station.StationUUID,
		}
	}
	return items
}

func (stationSource) category() category {
	return stations
}
