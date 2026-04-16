package world

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	"github.com/notkaj/grad/internal/selection/common"
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
		desc := fmt.Sprintf("Station %d", country.StationCount)
		items[i+1] = common.NewItem(country.Name, desc, country.Code)
		sum += country.StationCount
	}

	items[0] = common.AllItem
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
	order      g.StationsOrder
	by         g.StationsBy
	term       string
	reversed   bool
	hideBroken bool
}

func (s stationSource) items() []list.Item {
	var stations []g.Station
	if s.by == "" {
		stations = g.FetchAllStationsDetailed(s.order, s.reversed, 0, 100, s.hideBroken)
	} else {
		stations = g.FetchStationsDetailed(s.by, s.term, s.order, s.reversed, 0, 500, s.hideBroken)
	}
	len := len(stations)
	items := make([]list.Item, len)
	for i, station := range stations {
		var builder strings.Builder
		fmt.Fprintf(&builder, "%d clicks, ", station.ClickCount)
		fmt.Fprintf(&builder, "last okay check at %s", station.LastCheckOkTime)
		// if station.LastCheckOk {
		// 	fmt.Fprintf(&builder, "last check OK at %s", station.LastCheckOkTime)
		// } else {
		// 	fmt.Fprintf(&builder, "last check FAILED at %s", station.LastCheckTime)
		// }

		items[i] = common.NewItem(station.Name, builder.String(), station.StationUUID)
	}
	return items
}

func (stationSource) category() category {
	return stations
}

func defaultStationSource(countryCode string) stationSource {
	if countryCode == "ALL" {
		return allStationSource()
	}
	return stationSource{
		order:      g.StationsOrderClickCount,
		by:         g.StationsByCountryCodeExact,
		term:       countryCode,
		reversed:   true,
		hideBroken: true,
	}
}

func allStationSource() stationSource {
	return stationSource{
		order:      g.StationsOrderClickCount,
		reversed:   true,
		hideBroken: true,
	}
}
