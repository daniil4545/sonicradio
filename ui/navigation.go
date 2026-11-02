package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
	default: // up: already in the tab bar
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

func isTabBarKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "left", "right", "h", "l", "up", "down", "enter", "esc":
		return true
	}
	return false
}

func tabBarHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "tabs")),
		key.NewBinding(key.WithKeys("left", "right", "h", "l"), key.WithHelp("←/→ h/l", "switch tab in tab bar")),
		key.NewBinding(key.WithKeys("down", "enter", "esc"), key.WithHelp("↓/enter/esc", "return to list")),
	}
}
