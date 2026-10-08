package main

import (
	"errors"
	"strings"
	"testing"
)

// These tests guard the catalog against the mistakes that are easy to make
// when adding a tool: duplicates, missing install steps, missing "already
// installed" checks, and so on.

func TestCatalogIntegrity(t *testing.T) {
	cats := buildCategories()
	if len(cats) == 0 {
		t.Fatal("no categories defined")
	}

	seenKeys := map[string]bool{}
	seenTools := map[string]string{}

	for _, cat := range cats {
		if cat.Key == "" || cat.Title == "" {
			t.Errorf("category %+v needs a Key and a Title", cat)
		}
		if seenKeys[cat.Key] {
			t.Errorf("duplicate category key %q", cat.Key)
		}
		seenKeys[cat.Key] = true

		if len(cat.Tools) == 0 {
			t.Errorf("category %q has no tools", cat.Key)
		}

		for _, tool := range cat.Tools {
			name := strings.ToLower(tool.Name)
			if name == "" {
				t.Errorf("category %q contains a tool with no Name", cat.Key)
				continue
			}
			if prev, dup := seenTools[name]; dup {
				t.Errorf("tool %q appears in both %q and %q", tool.Name, prev, cat.Key)
			}
			seenTools[name] = cat.Key

			if tool.Method == "" {
				t.Errorf("%s: missing Method label (used by --list)", tool.Name)
			}
			if tool.Manual == "" && tool.Install == nil {
				t.Errorf("%s: needs an Install func or a Manual note", tool.Name)
			}
			if tool.Manual == "" && tool.Bin == "" && tool.CheckFunc == nil {
				t.Errorf("%s: needs Bin or CheckFunc so already-installed tools can be skipped", tool.Name)
			}
		}
	}
}

func TestPrerequisitesRunFirst(t *testing.T) {
	cats := buildCategories()
	if cats[0].Key != "prereq" {
		t.Fatalf("first category should be prereq (git/go/pipx/cargo are needed by everything after it), got %q", cats[0].Key)
	}
}

func TestCSVSet(t *testing.T) {
	got := csvSet(" Nuclei, httpx ,,SUBFINDER ")
	for _, want := range []string{"nuclei", "httpx", "subfinder"} {
		if !got[want] {
			t.Errorf("csvSet missing %q: %v", want, got)
		}
	}
	if len(got) != 3 {
		t.Errorf("expected 3 entries, got %d: %v", len(got), got)
	}
	if len(csvSet("")) != 0 {
		t.Error("empty input should yield an empty set")
	}
}

func TestSanitize(t *testing.T) {
	if got := sanitize("testssl.sh"); got != "testssl_sh" {
		t.Errorf("sanitize(testssl.sh) = %q", got)
	}
	if got := sanitize("binutils (strings/objdump)"); strings.ContainsAny(got, "/ \\") {
		t.Errorf("sanitize left path-unsafe characters in %q", got)
	}
}

func TestErrFromUsesLastNonEmptyLine(t *testing.T) {
	if errFrom("whatever", nil) != nil {
		t.Error("nil error must stay nil")
	}
	err := errFrom("first line\nreal problem here\n\n", errors.New("exit status 1"))
	if err == nil || err.Error() != "real problem here" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestShCapturesOutputAndStatus(t *testing.T) {
	out, err := sh("echo hello")
	if err != nil || strings.TrimSpace(out) != "hello" {
		t.Errorf("sh(echo hello) = %q, %v", out, err)
	}
	if _, err := sh("exit 3"); err == nil {
		t.Error("expected a non-zero exit to surface as an error")
	}
}
