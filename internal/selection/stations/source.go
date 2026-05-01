package stations

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/list"
	g "gitlab.com/notkaj/goradios"
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
	stations = g.FetchStationsDetailed(s.by, s.term, s.order, s.reversed, start, s.chunkSize, s.hideBroken)
	len := len(stations)
	items := make([]list.Item, len)
	for i, station := range stations {
		desc := fmt.Sprintf("%s, %d clicks", station.Codec, station.ClickCount)
		items[i] = NewItem(strings.TrimSpace(station.Name), desc, station.StationUUID, station.URLResolved, station.Homepage, station.Codec)
	}
	return items
}

func stationsByCountryCodeSource(countryCode string) source {
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

func (s *source) incrementChunk() {
	s.currentChunk += 1
}
