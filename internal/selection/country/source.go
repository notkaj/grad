package country

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	sel "github.com/notkaj/grad/internal/selection/selector"
	g "gitlab.com/AgentNemo/goradios"
)

type source struct {
	order      g.StationsOrder
	by         g.StationsBy
	term       string
	reversed   bool
	hideBroken bool
}

func (s source) items() []list.Item {
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

		items[i] = sel.NewItem(station.Name, builder.String(), station.StationUUID)
	}
	return items
}

func stationsByCountryCodeSource(countryCode string) source {
	if countryCode == "ALL" {
		return allStationSource()
	}
	return source{
		order:      g.StationsOrderClickCount,
		by:         g.StationsByCountryCodeExact,
		term:       countryCode,
		reversed:   true,
		hideBroken: true,
	}
}

func allStationSource() source {
	return source{
		order:      g.StationsOrderClickCount,
		reversed:   true,
		hideBroken: true,
	}
}
