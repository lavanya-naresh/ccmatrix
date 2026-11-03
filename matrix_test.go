package main

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
)

func TestMatrixRuns(t *testing.T) {
	cfg, _ := parseArgs([]string{"-c", "-r"}, io.Discard)
	m := newMatrix(30, 20, cfg, rand.New(rand.NewSource(1)))
	r := newRenderer(30, 20)
	var buf bytes.Buffer
	sawHead := false
	for i := 0; i < 2000; i++ {
		m.step()
		m.draw(r)
		for _, c := range r.next {
			if c.sgr == headSGR {
				sawHead = true
				if c.ch == ' ' || c.ch == 0 {
					t.Fatal("head cell has no character")
				}
			}
		}
		if err := r.flush(&buf); err != nil {
			t.Fatal(err)
		}
		buf.Reset()
	}
	if !sawHead {
		t.Error("never saw a head")
	}
}

func TestMatrixResize(t *testing.T) {
	cfg, _ := parseArgs(nil, io.Discard)
	m := newMatrix(10, 10, cfg, rand.New(rand.NewSource(2)))
	for _, sz := range [][2]int{{5, 5}, {1, 1}, {80, 24}, {2, 50}} {
		m.resize(sz[0], sz[1])
		r := newRenderer(sz[0], sz[1])
		for i := 0; i < 200; i++ {
			m.step()
			m.draw(r)
		}
	}
}

func TestRendererOnlyEmitsChanges(t *testing.T) {
	r := newRenderer(3, 2)
	r.next[0] = cell{'x', 32}
	var b bytes.Buffer
	if err := r.flush(&b); err != nil || b.Len() == 0 {
		t.Fatalf("first flush should draw, err=%v len=%d", err, b.Len())
	}
	b.Reset()
	if err := r.flush(&b); err != nil || b.Len() != 0 {
		t.Fatalf("second flush should be empty, err=%v len=%d", err, b.Len())
	}
}
