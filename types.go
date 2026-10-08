package main

import "os"

// Tool describes one installable thing.
type Tool struct {
	Name      string       // display name
	Method    string       // short label for --list: go/apt/pip/pipx/cargo/npm/git/script/manual
	Bin       string       // binary to look up on PATH to decide "already installed" (if CheckFunc is nil)
	CheckFunc func() bool  // custom "already installed" check; overrides Bin if set
	Install   func() error // how to install it; nil for Manual-only entries
	Manual    string       // if non-empty, this is a manual/how-to entry — Install is never called
	Info      string       // optional one-line note printed after a successful install
}

// Category groups related tools under a header.
type Category struct {
	Key   string
	Title string
	Icon  string
	Tools []Tool
}

func homeDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "/root"
	}
	return h
}
