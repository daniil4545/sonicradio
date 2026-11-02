package ui

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Retro header: tuner scale, fake spectrum analyzer and volume bar.

const (
	retroBars    = 12
	retroBarRows = 3
	tunerW       = 26 // scale columns for 88..108 MHz
	tunerMinFreq = 88.0
	tunerMaxFreq = 108.0
	// width below which the analyzer is hidden
	tunerGap              = 4
	retroMinWidthAnalyzer = tunerW + tunerGap + retroBars*2 + 2*HeaderPadDist
)

var retroBlocks = []rune(" ▁▂▃▄▅▆▇█")

// RetroColors are the extra colors of the retro header.
type RetroColors struct {
	text, dim, frame, label string
	led                     [retroBarRows]string // bottom, middle, top row
	// accents: tuner pointer, scale digits and frequency, song and tags
	needle, scale, song string
}

type retroState struct {
	frame       int
	level, peak []float64
	needle      float64
	rnd         *rand.Rand
}

func newRetroState() retroState {
	return retroState{
		level:  make([]float64, retroBars),
		peak:   make([]float64, retroBars),
		needle: tunerMinFreq,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec // visual noise only
	}
}

// stationFreq spreads n stations evenly over the FM scale.
func stationFreq(idx, n int) float64 {
	if n <= 0 || idx < 0 {
		return tunerMinFreq
	}
	return tunerMinFreq + (float64(idx)+0.5)*(tunerMaxFreq-tunerMinFreq)/float64(n)
}

func tunerCol(f float64) int {
	c := int(math.Round((f - tunerMinFreq) / (tunerMaxFreq - tunerMinFreq) * (tunerW - 1)))
	return max(0, min(tunerW-1, c))
}

// step advances one animation frame: the pointer glides to target,
// bars jump like music when playing, and show static while tuning.
func (r *retroState) step(target float64, playing bool) {
	r.frame++
	r.needle += (target - r.needle) * 0.3
	tuning := math.Abs(target-r.needle) >= 0.05
	if !tuning {
		r.needle = target
	}
	beat := r.frame%5 == 0
	for i := range r.level {
		var t float64
		switch {
		case !playing:
			t = 0
		case tuning:
			t = 0.05 + 0.3*r.rnd.Float64()
		default:
			env := 0.62 - 0.3*float64(i)/retroBars
			t = env * (0.3 + 0.7*r.rnd.Float64())
			if beat && i < 3 {
				t = math.Min(1, t+0.35)
			}
		}
		if tuning || !playing {
			r.level[i] = t
		} else {
			r.level[i] = math.Max(t, r.level[i]-0.1)
		}
		r.peak[i] = math.Max(r.level[i], r.peak[i]-0.03)
	}
}

func (r *retroState) tuning(target float64) bool { return math.Abs(target-r.needle) >= 0.05 }

func fg(c, s string) string { return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(s) }

func (r *retroState) tunerRows(c RetroColors) [retroBarRows]string {
	labels := []rune(strings.Repeat(" ", tunerW+tunerGap))
	for f := int(tunerMinFreq); f <= int(tunerMaxFreq); f += 4 {
		copy(labels[tunerCol(float64(f)):], []rune(strconv.Itoa(f)))
	}
	nc := tunerCol(r.needle)
	var scale, base strings.Builder
	for col := range tunerW {
		if col == nc {
			scale.WriteString(fg(c.needle, "┃"))
			base.WriteString(fg(c.needle, "▲"))
			continue
		}
		t := "╷"
		if col%5 == 0 {
			t = "│"
		}
		scale.WriteString(fg(c.dim, t))
		base.WriteString(fg(c.frame, "━"))
	}
	gap := strings.Repeat(" ", tunerGap)
	return [retroBarRows]string{fg(c.scale, string(labels)), scale.String() + gap, base.String() + gap}
}

func (r *retroState) analyzerRows(c RetroColors, playing bool) [retroBarRows]string {
	var rows [retroBarRows]string
	if !playing {
		blank := strings.Repeat(" ", retroBars*2)
		return [retroBarRows]string{blank, blank, blank}
	}
	for row := retroBarRows - 1; row >= 0; row-- {
		var b strings.Builder
		for i := range retroBars {
			fill := (r.level[i]*retroBarRows - float64(row)) * 8
			ch := retroBlocks[int(math.Max(0, math.Min(8, fill)))]
			pk := r.peak[i] * retroBarRows
			if ch == ' ' && playing && pk > float64(row) && pk <= float64(row+1) {
				ch = '▔'
			}
			b.WriteString(fg(c.led[row], string(ch)) + " ")
		}
		rows[retroBarRows-1-row] = b.String()
	}
	return rows
}

func volumeBarView(c RetroColors, vol int) string {
	vol = max(0, min(100, vol))
	n := vol / 5
	var on strings.Builder
	for i := range n { // low segments take the bottom LED color, loud ones the top
		on.WriteString(fg(c.led[min(retroBarRows-1, i*retroBarRows/20)], "▰"))
	}
	return fg(c.dim, "VOL ") + on.String() +
		fg(c.dim, strings.Repeat("▱", 20-n)) + fg(c.dim, fmt.Sprintf(" %3d", vol))
}

// retroTarget is the tuner frequency: the playing station, else the
// favorite under the cursor. found is false when the station is not a favorite.
func (m *Model) retroTarget() (freq float64, found bool) {
	favs := m.cfg.GetFavorites()
	m.delegate.playingMtx.RLock()
	s := m.delegate.currPlaying
	if s == nil {
		s = m.delegate.prevPlaying
	}
	m.delegate.playingMtx.RUnlock()
	if s != nil {
		for i := range favs {
			if favs[i].Stationuuid == s.Stationuuid {
				return stationFreq(i, len(favs)), true
			}
		}
		return m.retro.needle, false
	}
	if t, ok := m.tabs[favoriteTabIx].(stationTab); ok && t.Stations() != nil {
		return stationFreq(t.Stations().list.Index(), len(favs)), true
	}
	return tunerMinFreq, false
}

// retroStep is called on each spinner tick (10 fps while a station plays).
func (m *Model) retroStep() {
	if m.cfg.ClassicHeader || m.cfg.NoAnimation {
		return
	}
	target, _ := m.retroTarget()
	m.delegate.playingMtx.RLock()
	playing := m.delegate.currPlaying != nil
	m.delegate.playingMtx.RUnlock()
	m.retro.step(target, playing)
}

func (m *Model) retroHeaderView(width int) string {
	c := m.style.retro
	target, found := m.retroTarget()

	m.delegate.playingMtx.RLock()
	curr, prev := m.delegate.currPlaying, m.delegate.prevPlaying
	m.delegate.playingMtx.RUnlock()
	playing := curr != nil
	// no ticks run while nothing plays, so the pointer jumps instead of gliding
	if m.cfg.NoAnimation || !playing {
		m.retro.needle = target
	}

	var status string
	switch {
	case curr != nil && found && m.retro.tuning(target) && !m.cfg.NoAnimation:
		status = fg(c.label, "≈") + fg(c.dim, fmt.Sprintf(" tuning · %.1f FM", m.retro.needle))
	case curr != nil || prev != nil:
		s, icon := curr, "▶"
		if s == nil {
			s, icon = prev, "❚❚"
		}
		freq := ""
		if found {
			freq = fg(c.dim, " · ") + fg(c.scale, fmt.Sprintf("%.1f FM", target))
		}
		t := m.playbackTime
		status = fg(c.label, icon) + " " + fg(c.text, s.Name) + freq + fg(c.dim, " · ") +
			fg(c.text, fmt.Sprintf("%02d:%02d:%02d", int(t.Hours()), int(t.Minutes())%60, int(t.Seconds())%60))
	default:
		status = fg(c.dim, "■ "+noPlayingMsg)
	}
	if m.statusMsg != "" {
		status += fg(c.dim, "  "+m.statusMsg)
	}

	pad := strings.Repeat(" ", HeaderPadDist)
	var res strings.Builder
	res.WriteString(pad + fg(c.label, "SONIC DECK") + "  " +
		lipgloss.NewStyle().MaxWidth(max(0, width-2*HeaderPadDist-12)).Render(status) + "\n")

	tu := m.retro.tunerRows(c)
	an := m.retro.analyzerRows(c, playing)
	for i := range tu {
		line := pad + tu[i]
		if width >= retroMinWidthAnalyzer {
			line += an[i]
		}
		res.WriteString(line + "\n")
	}
	res.WriteString(pad + volumeBarView(c, m.cfg.GetVolume()) + "\n")

	song := ""
	if m.songTitle != "" {
		song = fg(c.needle, "♪ ") + fg(c.song, m.songTitle)
	}
	res.WriteString(pad + lipgloss.NewStyle().MaxWidth(max(0, width-2*HeaderPadDist)).Render(song) + "\n\n")
	return res.String()
}
