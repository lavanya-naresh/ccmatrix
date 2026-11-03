package main

import (
	"bytes"
	"io"
	"strconv"
)

// cell is one screen position. sgr is the SGR foreground code (0 = default).
type cell struct {
	ch  rune
	sgr uint8
}

var blankCell = cell{' ', 0}

// renderer double-buffers the screen and only emits cells that changed,
// batching everything into a single write per frame.
type renderer struct {
	cols, rows int
	prev, next []cell
	force      bool // redraw everything on next flush
	buf        bytes.Buffer
}

func newRenderer(cols, rows int) *renderer {
	r := &renderer{
		cols:  cols,
		rows:  rows,
		prev:  make([]cell, cols*rows),
		next:  make([]cell, cols*rows),
		force: true,
	}
	for i := range r.next {
		r.next[i] = blankCell
		r.prev[i] = blankCell
	}
	return r
}

func (r *renderer) set(col, row int, c cell) {
	r.next[row*r.cols+col] = c
}

func (r *renderer) flush(w io.Writer) error {
	r.buf.Reset()
	curRow, curCol, curSGR := -1, -1, -1
	for row := 0; row < r.rows; row++ {
		for col := 0; col < r.cols; col++ {
			i := row*r.cols + col
			n := r.next[i]
			if !r.force && n == r.prev[i] {
				continue
			}
			if curRow != row || curCol != col {
				r.buf.WriteString("\x1b[")
				r.buf.WriteString(strconv.Itoa(row + 1))
				r.buf.WriteByte(';')
				r.buf.WriteString(strconv.Itoa(col + 1))
				r.buf.WriteByte('H')
				curRow, curCol = row, col
			}
			if int(n.sgr) != curSGR {
				if n.sgr == 0 {
					r.buf.WriteString("\x1b[0m")
				} else {
					r.buf.WriteString("\x1b[")
					r.buf.WriteString(strconv.Itoa(int(n.sgr)))
					r.buf.WriteByte('m')
				}
				curSGR = int(n.sgr)
			}
			r.buf.WriteRune(n.ch)
			curCol++
			if curCol >= r.cols {
				curCol = -1 // cursor position after last column is unreliable
			}
		}
	}
	copy(r.prev, r.next)
	r.force = false
	if r.buf.Len() == 0 {
		return nil
	}
	r.buf.WriteString("\x1b[0m")
	_, err := w.Write(r.buf.Bytes())
	return err
}
