package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dancnb/sonicradio/model"
)

type dialTuneMsg struct {
	generation uint64
	uuid       string
}

// navigationList excludes forms, filters and station details.
func (m *Model) navigationList() *list.Model {
	switch t := m.tabs[m.activeTabIdx].(type) {
	case stationTab:
		if !t.IsFiltering() && !t.IsSearchEnabled() && !t.IsCustomStationEnabled() && !t.IsInfoEnabled() {
			return &t.Stations().list
		}
	case *historyTab:
		if !t.IsFiltering() {
			return &t.list
		}
	}
	return nil
}

func (m *Model) dialStation() string {
	if m.activeTabIdx != favoriteTabIx || m.tabBarFocused {
		return ""
	}
	l := m.navigationList()
	if l == nil {
		return ""
	}
	s, ok := l.SelectedItem().(model.Station)
	m.delegate.playingMtx.RLock()
	defer m.delegate.playingMtx.RUnlock()
	if !ok || m.delegate.currPlaying == nil || s.Stationuuid == m.delegate.currPlaying.Stationuuid {
		return ""
	}
	return s.Stationuuid
}

func (m *Model) updateTabBar(msg tea.KeyMsg) tea.Cmd {
	next := int(m.activeTabIdx)
	switch msg.String() {
	case "left", "h":
		next = (next + len(m.tabs) - 1) % len(m.tabs)
	case "right", "l":
		next = (next + 1) % len(m.tabs)
	case "down", "enter", "esc":
		m.tabBarFocused = false
		return nil
	default:
		return nil
	}
	var exit tea.Cmd
	if m.activeTabIdx == settingsTabIx {
		if settings, ok := m.tabs[settingsTabIx].(*settingsTab); ok {
			exit = settings.onExit(m)
		}
	}
	switch next {
	case int(favoriteTabIx):
		m.toFavoritesTab()
	case int(browseTabIx):
		m.toBrowseTab()
	case int(historyTabIx):
		m.toHistoryTab()
	case int(settingsTabIx):
		return m.toSettingsTab()
	}
	return exit
}

func tabBarHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "tabs")),
		key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→ h/l", "switch tab in tab bar")),
		key.NewBinding(key.WithKeys("down", "enter", "esc"), key.WithHelp("↓/enter/esc", "return to list")),
	}
}
