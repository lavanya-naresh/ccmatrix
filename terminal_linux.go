//go:build linux

package main

import "syscall"

const (
	ioctlGetTermios      = syscall.TCGETS
	ioctlSetTermios      = syscall.TCSETS
	ioctlSetTermiosFlush = 0x5404
	// ioctlSetTermiosFlush = syscall.TCSETSF
)
