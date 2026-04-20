package country

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	g "gitlab.com/AgentNemo/goradios"
)

type source struct {
	order        g.StationsOrder
	by           g.StationsBy
	term         string
	reversed     bool
	hideBroken   bool
	currentChunk uint
	chunkSize    uint
}

func (s source) items() []list.Item {
	var stations []g.Station
	start := s.currentChunk * s.chunkSize
	if s.by == "" {
		stations = g.FetchAllStationsDetailed(s.order, s.reversed, start, s.chunkSize, s.hideBroken)
	} else {
		stations = g.FetchStationsDetailed(s.by, s.term, s.order, s.reversed, start, s.chunkSize, s.hideBroken)
	}
	len := len(stations)
	items := make([]list.Item, len)
	for i, station := range stations {
		// var builder strings.Builder
		// fmt.Fprintf(&builder, "%d clicks, ", station.ClickCount)
		// fmt.Fprintf(&builder, "last okay check at %s", station.LastCheckOkTime)
		// if station.LastCheckOk {
		// 	fmt.Fprintf(&builder, "last check OK at %s", station.LastCheckOkTime)
		// } else {
		// 	fmt.Fprintf(&builder, "last check FAILED at %s", station.LastCheckTime)
		// }
		desc := fmt.Sprintf("%d clicks", station.ClickCount)
		items[i] = NewItem(station.Name, desc, station.StationUUID, station.URLResolved)
	}
	return items
}

func stationsByCountryCodeSource(countryCode string) source {
	if countryCode == "ALL" {
		return allStationSource()
	}
	return source{
		order:        g.StationsOrderClickCount,
		by:           g.StationsByCountryCodeExact,
		term:         countryCode,
		reversed:     true,
		hideBroken:   true,
		currentChunk: 0,
		chunkSize:    50,
	}
}

func allStationSource() source {
	return source{
		order:        g.StationsOrderClickCount,
		reversed:     true,
		hideBroken:   true,
		currentChunk: 0,
		chunkSize:    50,
	}
}

func (s *source) incrementChunk() {
	s.currentChunk += 1
}
