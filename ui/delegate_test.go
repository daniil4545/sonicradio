package ui

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/dancnb/sonicradio/config"
	"github.com/dancnb/sonicradio/model"
)

func TestMinimalStationLine(t *testing.T) {
	for _, tt := range []struct {
		name, tags string
		width      int
		want       string
	}{
		{"Radio", "jazz", 40, "Radio  jazz"},
		{"Radio", "", 40, "Radio"},
		{"Radio", "jazz", 8, "Radio  j"},
		{"Radio", "jazz", 5, "Radio"},
		{"Радио", "джаз", 8, "Радио  д"},
		{"世界電台", "jazz", 6, "世界電"},
	} {
		d := newStationDelegate(&config.Value{}, NewStyle(0), nil, nil)
		s := model.Station{Stationuuid: "station", Name: tt.name, Tags: tt.tags, Bitrate: 128, Votes: 10}
		prefix := d.style.PrefixStyle.Render(IndexString(1))
		width := tt.width + lipgloss.Width(prefix) + HeaderPadDist
		m := list.New([]list.Item{s}, d, width, 10)
		var out bytes.Buffer
		d.Render(&out, m, 0, s)
		got := ansi.Strip(out.String())
		if !utf8.ValidString(got) || strings.Contains(got, "\n") || lipgloss.Width(got) > width {
			t.Errorf("invalid station line: %q", got)
		}
		if got = strings.TrimSpace(strings.TrimPrefix(got, ansi.Strip(prefix))); got != tt.want {
			t.Errorf("line at width %d = %q, want %q", tt.width, got, tt.want)
		}
		t.Logf("station line: %q", got)
	}
}

func TestStationShortHelp(t *testing.T) {
	d := newStationDelegate(&config.Value{}, NewStyle(0), nil, nil)
	help := d.ShortHelp()
	want := []string{"enter/l", "space", "-", "+", "f"}
	if len(help) != len(want) {
		t.Fatalf("short help has %d bindings, want %d", len(help), len(want))
	}
	for i, binding := range help {
		if got := binding.Help().Key; got != want[i] {
			t.Errorf("short help key %d = %q, want %q", i, got, want[i])
		}
	}
}

func TestStationViewCycle(t *testing.T) {
	cfg := &config.Value{}
	d := newStationDelegate(cfg, NewStyle(0), nil, nil)
	m := Model{cfg: cfg, delegate: d}
	for _, want := range []config.StationView{config.MinimalView, config.CompactView, config.DefaultView, config.MinimalView} {
		if cfg.StationView != want {
			t.Fatalf("view = %s, want %s", cfg.StationView, want)
		}
		if want == config.MinimalView && (d.Height() != 1 || d.Spacing() != 0) {
			t.Fatal("minimal view must use one line without spacing")
		}
		m.changeStationView()
	}
}
