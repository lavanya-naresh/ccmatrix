package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// os.Exit skips defers, so all real work (and terminal restore) lives in run.
	os.Exit(run(os.Args[1:]))
}

// intervalFor maps the -u level (0..10) to a tick interval, with a floor so
// level 0 is fast but not a busy loop.
func intervalFor(level int) time.Duration {
	d := time.Duration(level) * 10 * time.Millisecond
	if d < 5*time.Millisecond {
		d = 5 * time.Millisecond
	}
	return d
}

func run(args []string) (code int) {
	cfg, err := parseArgs(args, os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		fmt.Fprintf(os.Stderr, "ccmatrix: %v\n", err)
		fmt.Fprintln(os.Stderr, "Try 'ccmatrix -h' for usage.")
		return 2
	}

	if !isTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "ccmatrix: standard input must be a terminal")
		return 1
	}

	term, err := newTerminal(cfg.Screensaver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccmatrix: cannot set up terminal: %v\n", err)
		return 1
	}
	defer term.Restore()
	defer func() {
		if r := recover(); r != nil {
			term.Restore()
			fmt.Fprintf(os.Stderr, "ccmatrix: internal error: %v\n", r)
			code = 1
		}
	}()

	sigs := make(chan os.Signal, 4)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)

	keys := make(chan keyEvent, 16)
	go readInput(os.Stdin, keys)

	out := os.Stdout
	cols, rows := term.Size()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := newMatrix(cols, rows, cfg, rng)
	r := newRenderer(cols, rows)

	level := cfg.Speed
	ticker := time.NewTicker(intervalFor(level))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.step()
			m.draw(r)
			if err := r.flush(out); err != nil {
				return 1
			}

		case ev, ok := <-keys:
			if !ok || cfg.Screensaver {
				return 0
			}
			if ev.esc {
				continue // arrows, function keys etc. are ignored
			}
			switch {
			case ev.b == 'q' || ev.b == 'Q' || ev.b == 3: // 3 = Ctrl-C in raw mode
				return 0
			case ev.b >= '0' && ev.b <= '9':
				level = int(ev.b - '0')
				ticker.Reset(intervalFor(level))
			}

		case <-winch:
			nc, nr := term.Size()
			if nc != cols || nr != rows {
				cols, rows = nc, nr
				m.resize(cols, rows)
				r = newRenderer(cols, rows) // force a full redraw
				out.WriteString("\x1b[2J")
			}

		case <-sigs:
			return 0
		}
	}
}
