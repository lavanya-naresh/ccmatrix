package main

import (
	"io"
	"testing"
)

func TestParseArgsValid(t *testing.T) {
	for _, args := range [][]string{
		{}, {"-s"}, {"-u", "0"}, {"-u", "10"}, {"-C", "red", "-c", "-r"}, {"-C", "Cyan", "-u", "2"},
	} {
		if _, err := parseArgs(args, io.Discard); err != nil {
			t.Errorf("args %v: unexpected error %v", args, err)
		}
	}
}

func TestParseArgsInvalid(t *testing.T) {
	for _, args := range [][]string{
		{"-u", "11"}, {"-u", "-1"}, {"-u", "fast"}, {"-C", "purple"}, {"extra"},
	} {
		if _, err := parseArgs(args, io.Discard); err == nil {
			t.Errorf("args %v: expected error", args)
		}
	}
}

func TestParseArgsDefaults(t *testing.T) {
	cfg, err := parseArgs(nil, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Speed != 4 || cfg.Colour != "green" || cfg.ColourCode != 32 || cfg.Screensaver || cfg.Rainbow || cfg.Katakana {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}
