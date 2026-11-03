package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

// SGR foreground colour codes for the supported trail colours.
var colourCodes = map[string]int{
	"green": 32, "red": 31, "blue": 34, "white": 37,
	"yellow": 33, "cyan": 36, "magenta": 35, "black": 30,
}

// Order used when listing colours in error messages.
var colourOrder = []string{"green", "red", "blue", "white", "yellow", "cyan", "magenta", "black"}

// Config holds every option the user can set on the command line.
type Config struct {
	Screensaver bool   // -s
	Speed       int    // -u (0 fastest .. 10 slowest)
	Colour      string // -C
	ColourCode  int    // SGR base code resolved from Colour
	Rainbow     bool   // -r
	Katakana    bool   // -c
}

// parseArgs parses and validates args. It never touches the terminal, so
// errors are reported before anything is drawn. It returns flag.ErrHelp for -h.
func parseArgs(args []string, stderr io.Writer) (Config, error) {
	var cfg Config
	fs := flag.NewFlagSet("ccmatrix", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // we print our own errors and usage
	fs.BoolVar(&cfg.Screensaver, "s", false, "")
	fs.IntVar(&cfg.Speed, "u", 4, "")
	fs.StringVar(&cfg.Colour, "C", "green", "")
	fs.BoolVar(&cfg.Rainbow, "r", false, "")
	fs.BoolVar(&cfg.Katakana, "c", false, "")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printUsage(stderr)
		}
		return cfg, err
	}
	if fs.NArg() > 0 {
		return cfg, fmt.Errorf("unexpected argument %q", fs.Arg(0))
	}
	if cfg.Speed < 0 || cfg.Speed > 10 {
		return cfg, fmt.Errorf("-u must be a number from 0 to 10, got %d", cfg.Speed)
	}
	cfg.Colour = strings.ToLower(cfg.Colour)
	code, ok := colourCodes[cfg.Colour]
	if !ok {
		return cfg, fmt.Errorf("unknown colour %q; supported colours: %s",
			cfg.Colour, strings.Join(colourOrder, ", "))
	}
	cfg.ColourCode = code
	return cfg, nil
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `Usage: ccmatrix [options]

Options:
  -s         screensaver mode: exit on any key press or mouse activity
  -u N       delay between updates, 0 (fastest) to 10 (slowest), default 4
             (press 0-9 while running to change it)
  -C COLOUR  trail colour: green (default), red, blue, white, yellow,
             cyan, magenta, black
  -r         rainbow-coloured trails
  -c         use half-width katakana instead of letters and numbers
  -h         show this help

Press q to quit (unless -s is given).
`)
}
