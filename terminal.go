package main

import (
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type winsize struct{ Row, Col, Xpixel, Ypixel uint16 }

func ioctl(fd, req uintptr, arg unsafe.Pointer) error {
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(arg)); e != 0 {
		return e
	}
	return nil
}

func isTerminal(fd int) bool {
	var t syscall.Termios
	return ioctl(uintptr(fd), ioctlGetTermios, unsafe.Pointer(&t)) == nil
}

// Terminal owns raw mode, the alternate screen and (optionally) mouse
// reporting, and puts everything back in Restore.
type Terminal struct {
	fd    int
	old   syscall.Termios
	mouse bool
	once  sync.Once
}

func newTerminal(mouse bool) (*Terminal, error) {
	t := &Terminal{fd: int(os.Stdin.Fd()), mouse: mouse}
	if err := ioctl(uintptr(t.fd), ioctlGetTermios, unsafe.Pointer(&t.old)); err != nil {
		return nil, err
	}

	raw := t.old
	raw.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP |
		syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0
	if err := ioctl(uintptr(t.fd), ioctlSetTermios, unsafe.Pointer(&raw)); err != nil {
		return nil, err
	}

	// alt screen, hide cursor, disable line wrap, clear
	seq := "\x1b[?1049h\x1b[?25l\x1b[?7l\x1b[2J"
	if mouse {
		// any-event tracking + SGR encoding: clicks and movement arrive as stdin bytes
		seq += "\x1b[?1003h\x1b[?1006h"
	}
	os.Stdout.WriteString(seq)
	return t, nil
}

// Restore is idempotent and safe to call from defer, signal and panic paths.
func (t *Terminal) Restore() {
	t.once.Do(func() {
		seq := "\x1b[0m"
		if t.mouse {
			seq += "\x1b[?1006l\x1b[?1003l"
		}
		seq += "\x1b[?7h\x1b[?25h\x1b[?1049l"
		os.Stdout.WriteString(seq)
		// flush pending input so stray mouse reports don't leak into the shell
		ioctl(uintptr(t.fd), ioctlSetTermiosFlush, unsafe.Pointer(&t.old))
	})
}

// Size returns the terminal size in columns and rows (80x24 if unknown).
func (t *Terminal) Size() (cols, rows int) {
	var ws winsize
	if err := ioctl(uintptr(t.fd), syscall.TIOCGWINSZ, unsafe.Pointer(&ws)); err != nil ||
		ws.Col == 0 || ws.Row == 0 {
		return 80, 24
	}
	return int(ws.Col), int(ws.Row)
}
