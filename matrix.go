package main

import "math/rand"

const headSGR = 97 // bright white

var rainbowPalette = []uint8{31, 33, 32, 36, 34, 35}

// column is the state of one falling trail.
type column struct {
	head, length int // head row (-1 = not on screen yet), trail length
	period, tick int // advance one row every `period` ticks
	idle         int // ticks to wait before (re)starting
	hueOffset    int // rainbow phase
}

type matrix struct {
	cols, rows int
	cfg        Config
	rng        *rand.Rand
	charset    []rune
	columns    []column
	chars      []rune // chars[col*rows+row]; written once as the head passes
}

func newMatrix(cols, rows int, cfg Config, rng *rand.Rand) *matrix {
	m := &matrix{cfg: cfg, rng: rng, charset: asciiChars}
	if cfg.Katakana {
		m.charset = katakanaChars
	}
	m.resize(cols, rows)
	return m
}

// resize discards all state and starts fresh at the new size.
func (m *matrix) resize(cols, rows int) {
	m.cols, m.rows = cols, rows
	m.columns = make([]column, cols)
	m.chars = make([]rune, cols*rows)
	for i := range m.columns {
		m.respawn(&m.columns[i], true)
	}
}

func (m *matrix) respawn(c *column, initial bool) {
	minLen := max(3, m.rows/4)
	maxLen := max(minLen+1, m.rows*3/4)
	c.length = minLen + m.rng.Intn(maxLen-minLen)
	c.period = 1 + m.rng.Intn(3)
	c.tick = 0
	c.head = -1
	c.hueOffset = m.rng.Intn(len(rainbowPalette))
	if initial {
		c.idle = m.rng.Intn(m.rows + 1)
	} else {
		c.idle = m.rng.Intn(m.rows/2 + 1)
	}
}

func (m *matrix) randomChar() rune {
	return m.charset[m.rng.Intn(len(m.charset))]
}

// step advances the simulation by one tick.
func (m *matrix) step() {
	for col := range m.columns {
		c := &m.columns[col]
		if c.idle > 0 {
			c.idle--
			continue
		}
		c.tick++
		if c.tick < c.period {
			continue
		}
		c.tick = 0
		c.head++
		if c.head >= 0 && c.head < m.rows {
			m.chars[col*m.rows+c.head] = m.randomChar()
		}
		if c.head-c.length >= m.rows {
			m.respawn(c, false)
		}
	}
}

// draw fills the renderer's next buffer. Style is derived from the distance
// to the head, so the tail clears itself without any bookkeeping.
func (m *matrix) draw(r *renderer) {
	for col := 0; col < m.cols; col++ {
		c := &m.columns[col]
		for row := 0; row < m.rows; row++ {
			d := c.head - row
			if d < 0 || d > c.length {
				r.set(col, row, blankCell)
				continue
			}
			ch := m.chars[col*m.rows+row]
			if ch == 0 {
				r.set(col, row, blankCell)
				continue
			}
			if d == 0 {
				r.set(col, row, cell{ch, headSGR})
				continue
			}
			r.set(col, row, cell{ch, m.trailColour(c, d, row)})
		}
	}
}

func (m *matrix) trailColour(c *column, d, row int) uint8 {
	base := uint8(m.cfg.ColourCode)
	if m.cfg.Rainbow {
		base = rainbowPalette[(row+c.hueOffset)%len(rainbowPalette)]
	}
	if d <= c.length/3 {
		return base + 60 // bright variant near the head
	}
	return base
}
