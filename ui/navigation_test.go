package ui

import (
	"context"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dancnb/sonicradio/config"
	"github.com/dancnb/sonicradio/model"
)

func dialModel() *Model {
	cfg := &config.Value{}
	historyMax := config.DefHistorySaveMax
	cfg.HistorySaveMax = &historyMax
	stations := []model.Station{{Stationuuid: "a", Name: "A"}, {Stationuuid: "b", Name: "B"}, {Stationuuid: "c", Name: "C"}}
	cfg.SetFavorites(stations)
	s := NewStyle(0)
	d := newStationDelegate(cfg, s, nil, nil)
	f := &favoritesTab{stationsTabBase: newStationsTab(newListKeymap(), nil, s), cfg: cfg}
	f.list = f.createList(d, 100, 30)
	f.list.SetItems([]list.Item{stations[0], stations[1], stations[2]})
	b := &browseTab{stationsTabBase: newStationsTab(newListKeymap(), nil, s), searchModel: &searchModel{}}
	b.list = b.createList(d, 100, 30)
	b.list.SetItems(f.list.Items())
	h := &historyTab{style: s}
	h.createList(100, 30)
	h.list.SetItems([]list.Item{config.HistoryEntry{Uuid: "a"}, config.HistoryEntry{Uuid: "b"}})
	d.currPlaying = &stations[0]
	settings := newSettingsTab(context.Background(), cfg, s, []config.PlayerType{config.Internal}, func(int) {})
	return &Model{
		cfg: cfg, style: s, delegate: d, retro: newRetroState(), width: 100,
		tabs: []uiTab{f, b, h, settings}, statusUpdate: make(chan struct{}, 10),
	}
}

func press(m *Model, k tea.KeyType) tea.Cmd {
	_, cmd := m.Update(tea.KeyMsg{Type: k})
	return cmd
}

func TestPlayKeys(t *testing.T) {
	space := tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}
	enter := tea.KeyMsg{Type: tea.KeyEnter}
	for _, tc := range []struct {
		name    string
		key     tea.KeyMsg
		down    bool // cursor on the second station, A plays
		paused  bool
		connect bool // expect a new station to start
	}{
		{"space on playing pauses", space, false, false, false},
		{"enter on playing pauses", enter, false, false, false},
		{"space on other plays it", space, true, false, true},
		{"enter on other plays it", enter, true, false, true},
		{"enter on paused resumes", enter, false, true, false},
		{"space on other while paused plays it", space, true, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := dialModel()
			if tc.paused {
				m.delegate.prevPlaying, m.delegate.currPlaying = m.delegate.currPlaying, nil
			}
			if tc.down {
				press(m, tea.KeyDown)
			}
			_, cmd := m.Update(tc.key)
			if cmd == nil {
				t.Fatal("key did nothing")
			}
			if got := strings.HasPrefix(m.statusMsg, "Connecting to"); got != tc.connect {
				t.Fatalf("connect = %v, status %q", got, m.statusMsg)
			}
		})
	}
}

func TestCursorDoesNotSwitch(t *testing.T) {
	m := dialModel()
	press(m, tea.KeyDown)
	press(m, tea.KeyDown)
	if m.statusMsg != "" || m.delegate.currPlaying.Stationuuid != "a" {
		t.Fatal("moving the cursor must not switch stations")
	}
	if freq, _ := m.retroTarget(); freq != stationFreq(0, 3) {
		t.Fatal("pointer must stay on the playing station")
	}
}

func TestTabFocus(t *testing.T) {
	for _, tab := range []uiTabIndex{favoriteTabIx, browseTabIx, historyTabIx} {
		for _, exit := range []tea.KeyType{tea.KeyDown, tea.KeyEnter, tea.KeyEsc} {
			m := dialModel()
			m.activeTabIdx = tab
			press(m, tea.KeyUp)
			if !m.tabBarFocused || m.navigationList().Index() != 0 {
				t.Fatal("Up on first row did not focus tabs")
			}
			if !strings.Contains(m.headerView(100), "▸"+strings.TrimSpace(tab.String())) {
				t.Fatal("focused tab marker missing")
			}
			press(m, tea.KeyLeft)
			press(m, tea.KeyRight)
			m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
			m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
			if m.activeTabIdx != tab || !m.tabBarFocused {
				t.Fatal("tab switching lost focus or selected wrong tab")
			}
			press(m, exit)
			if m.tabBarFocused || m.navigationList().Index() != 0 {
				t.Fatal("return to list moved cursor")
			}
			l := m.navigationList()
			l.Select(len(l.Items()) - 1)
			press(m, tea.KeyDown)
			if l.Index() != 0 || m.tabBarFocused {
				t.Fatal("Down on last row must wrap")
			}
		}
	}
}

func TestDialHeaderAndSettings(t *testing.T) {
	m := dialModel()
	m.delegate.currPlaying = nil
	press(m, tea.KeyDown)
	start := m.retro.needle
	m.retroStep()
	if m.retro.needle <= start || m.retro.needle >= stationFreq(1, 3) {
		t.Fatal("animated pointer must glide toward the cursor")
	}
	m.cfg.NoAnimation = true
	m.headerView(100)
	if m.retro.needle != stationFreq(1, 3) {
		t.Fatal("NoAnimation pointer must jump to the cursor")
	}
	press(m, tea.KeyUp)
	press(m, tea.KeyUp)
	press(m, tea.KeyLeft)
	if m.activeTabIdx != settingsTabIx || !m.tabBarFocused {
		t.Fatal("tab bar must reach Settings")
	}
	press(m, tea.KeyDown)
	press(m, tea.KeyUp)
	if m.tabBarFocused {
		t.Fatal("Settings Up must retain form navigation")
	}
	f := m.tabs[favoriteTabIx].(*favoritesTab)
	f.list.Help.Width = 160
	help := f.list.Help.View(f.list)
	if !strings.Contains(help, "↑ tabs") {
		t.Fatal("wide short help must show tabs")
	}
	f.list.Help.ShowAll = true
	if !strings.Contains(f.list.Help.View(f.list), "return to list") {
		t.Fatal("full help must show tab bar controls")
	}
	t.Log("Favorites cursor: 98.0 FM; NoAnimation jumps; Settings keeps Up/Down; short help: ↑ tabs")
}

func TestTabBarPassesOtherKeys(t *testing.T) {
	m := dialModel()
	press(m, tea.KeyUp)
	if !m.tabBarFocused {
		t.Fatal("Up on first row did not focus tabs")
	}
	if !isTabBarKey(tea.KeyMsg{Type: tea.KeyRight}) || isTabBarKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}) {
		t.Fatal("tab bar key set is wrong")
	}
}
