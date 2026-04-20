// Package country provides selector for stations
package country

import (
	"time"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	sel "github.com/notkaj/grad/internal/selection/selector"
	s "github.com/notkaj/grad/internal/style"
)

type searchDebounceMsg struct {
	term string
	seq  int
}

type searchResultMsg []list.Item

type Model struct {
	list       list.Model
	source     source
	baseSource source
	baseTitle  string
	fetching   bool
	totalCount int
	inSearch   bool
	searchSeq  int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		s.LightDark = lipgloss.LightDark(msg.IsDark())
		m.list.Styles.Title = s.Styles.Title
		return m, nil

	case tea.WindowSizeMsg:
		width, height := s.Styles.App.GetFrameSize()
		m.list.SetSize(msg.Width-width, msg.Height-height)
		return m, nil

	case sel.PopulatedMsg:
		m.list.StopSpinner()
		m.fetching = false
		m.source.incrementChunk()
		return m, m.list.SetItems(append(m.list.Items(), msg...))

	case searchDebounceMsg:
		if msg.seq == m.searchSeq && msg.term != "" {
			cmds = append(cmds, m.list.StartSpinner(), searchFetch(msg.term))
		}
		return m, tea.Batch(cmds...)

	case searchResultMsg:
		m.list.StopSpinner()
		return m, m.list.SetItems([]list.Item(msg))
	}

	prevTerm := m.list.FilterInput.Value()

	newListModel, cmd := m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.list = newListModel

	newTerm := m.list.FilterInput.Value()
	newState := m.list.FilterState()

	switch {
	case m.inSearch && (newState == list.Unfiltered || newTerm == ""):
		m.inSearch = false
		m.source = m.baseSource
		m.source.currentChunk = 0
		m.list.Title = m.baseTitle
		m.totalCount = 0
		m.fetching = false
		var empty []list.Item
		m.list.SetItems(empty)
		cmds = append(cmds, m.Populate())

	case newTerm != prevTerm && newTerm != "":
		if !m.inSearch {
			m.inSearch = true
		}
		m.searchSeq++
		seq := m.searchSeq
		term := newTerm
		cmds = append(cmds, tea.Tick(300*time.Millisecond, func(time.Time) tea.Msg {
			return searchDebounceMsg{term: term, seq: seq}
		}))

	case !m.inSearch &&
		m.list.Paginator.OnLastPage() &&
		!m.fetching &&
		!m.IsFiltering() &&
		len(m.list.Items()) < m.totalCount:
		m.fetching = true
		cmds = append(cmds, m.Populate())
	}

	return m, tea.Batch(cmds...)
}

func searchFetch(term string) tea.Cmd {
	return func() tea.Msg {
		source := searchSource(term)
		return searchResultMsg(source.items())
	}
}

func (m *Model) View() tea.View {
	return tea.NewView(s.Styles.App.Render(m.list.View()))
}

func (m Model) Info() string {
	i, ok := m.list.SelectedItem().(item)
	if ok {
		return i.id
	}
	return ""
}

func (m *Model) Populate() tea.Cmd {
	m.fetching = true
	return tea.Batch(
		m.list.StartSpinner(),
		func() tea.Msg {
			return sel.PopulatedMsg(m.source.items())
		},
	)
}

func (m *Model) Select(msg sel.Msg) {
	m.source.currentChunk = 0
	m.inSearch = false
	m.searchSeq++
	m.list.ResetFilter()
	var items []list.Item
	m.list.SetItems(items)
	switch msg := msg.(type) {
	case sel.CountrySelectedMsg:
		code := msg.Code
		size := msg.Count
		m.totalCount = size
		m.list.Title = msg.Name
		if code == "ALL" {
			m.source = allStationSource()
		} else {
			m.source = stationsByCountryCodeSource(code)
		}
		m.baseSource = m.source
		m.baseTitle = m.list.Title
	}
}

func (m *Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

func (m *Model) ViewLayer() *lipgloss.Layer {
	return lipgloss.NewLayer(s.Styles.App.Render(m.list.View()))
}

func (m *Model) SelectionInfo() (string, string, string) {
	if i, ok := m.list.SelectedItem().(item); ok {
		return i.title, i.url, i.codec
	}
	return "", "", ""
}

func InitialModel() Model {
	m := Model{}
	s.LightDark = lipgloss.LightDark(true)

	delegate := newItemDelegate()
	delegate.Styles.SelectedTitle = s.Styles.SeletectedTitle
	delegate.Styles.SelectedDesc = s.Styles.SelectedDesc
	stationsList := list.New(nil, delegate, 0, 0)
	stationsList.Title = "Stations"
	stationsList.Styles.Title = s.Styles.Title
	stationsList.Filter = func(_ string, targets []string) []list.Rank {
		ranks := make([]list.Rank, len(targets))
		for i := range targets {
			ranks[i] = list.Rank{Index: i}
		}
		return ranks
	}

	m.list = stationsList
	m.fetching = false
	m.list.SetSpinner(spinner.Line)

	return m
}
