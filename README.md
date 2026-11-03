# ccmatrix

Matrix-style falling text for your terminal (Coding Challenge #144), written in Go
with **no third-party dependencies** - just the standard library, `syscall` for
raw mode, and ANSI escape codes. Works on Linux and macOS.

## Build & run

Requires Go 1.21+.

    go build -o ccmatrix .
    ./ccmatrix

Cross-compile for all four targets (linux/darwin x amd64/arm64) into `dist/`:

    ./build.sh

Run the tests: `go test ./...`

## Options

    -s         screensaver mode: exit on any key press or mouse activity
    -u N       delay between updates, 0 (fastest) to 10 (slowest), default 4
               (press 0-9 while running to change it)
    -C COLOUR  green (default), red, blue, white, yellow, cyan, magenta, black
    -r         rainbow trails
    -c         half-width katakana instead of letters and numbers
    -h         help

Press `q` (or Ctrl-C) to quit. Examples:

    ./ccmatrix -c -C cyan -u 2
    ./ccmatrix -s -r

## Layout

    main.go      event loop: ticker, keys, resize, signals
    config.go    flag parsing and validation (before the terminal is touched)
    terminal.go  raw mode, alt screen, mouse, size, idempotent restore
    terminal_{linux,darwin}.go   per-OS termios ioctl numbers
    input.go     stdin reader goroutine; escape sequences collapse to one event
    matrix.go    simulation: per-column heads, speeds and trail lengths
    render.go    double-buffered diff renderer, one write per frame
