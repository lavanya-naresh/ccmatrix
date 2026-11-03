package main

import "io"

// keyEvent is one key press. Escape sequences (arrows, function keys, mouse
// reports) are collapsed into a single event with esc=true so their inner
// digits and letters are never mistaken for 'q' or speed keys.
type keyEvent struct {
	b   byte
	esc bool
}

// readInput forwards stdin as key events and closes out on EOF/error.
func readInput(r io.Reader, out chan<- keyEvent) {
	defer close(out)
	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			if buf[0] == 0x1b {
				out <- keyEvent{esc: true}
			} else {
				for _, b := range buf[:n] {
					out <- keyEvent{b: b}
				}
			}
		}
		if err != nil {
			return
		}
	}
}
