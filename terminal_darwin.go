//go:build darwin

package main

import "syscall"

const (
	ioctlGetTermios      = syscall.TIOCGETA
	ioctlSetTermios      = syscall.TIOCSETA
	ioctlSetTermiosFlush = syscall.TIOCSETAF
)
