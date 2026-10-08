// bb-arsenal — Bug Bounty Arsenal Installer
//
// A single-binary, dependency-free Go program that installs a curated set of
// recon / bug-bounty CLI tools on Debian/Ubuntu/Kali-style Linux systems,
// skipping anything already installed, with an animated terminal UI.
//
// Build:   go build -o bb-arsenal .
// Run:     ./bb-arsenal            (installs everything, asks first)
//
//	./bb-arsenal --list     (show the full catalog, install nothing)
//	./bb-arsenal -y         (skip the confirmation prompt)
//	./bb-arsenal --dry-run  (show what would happen)
//	./bb-arsenal --categories subs,dns,vuln
//	./bb-arsenal --only nuclei,httpx,subfinder
//	./bb-arsenal --skip masscan,rustscan
//	./bb-arsenal --no-gui   (skip the optional GUI-app category)
//	./bb-arsenal --version
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────
// Styling
// ─────────────────────────────────────────────────────────────────────────

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[90m"
	white   = "\033[97m"
)

var noColor bool

// version is overridden at release time: -ldflags "-X main.version=v1.0.0"
var version = "dev"

const repoURL = "https://github.com/BlackRootSec/bb-arsenal"

func c(color, s string) string {
	if noColor {
		return s
	}
	return color + s + reset
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const divider = "──────────────────────────────────────────────────────────────"

func banner() {
	art := `
 ____             ____                   _         _
| __ ) _   _  __ _| __ )  ___  _   _ _ __ | |_ _   _| |
|  _ \| | | |/ _  |  _ \ / _ \| | | |  _ \| __| | | | |
| |_) | |_| | (_| | |_) | (_) | |_| | | | | |_| |_| |_|
|____/ \__,_|\__, |____/ \___/ \__,_|_| |_|\__|\__, (_)
             |___/                             |___/
              A R S E N A L   I N S T A L L E R
`
	fmt.Println(c(cyan+bold, art))
}

func sectionHeader(icon, title string) {
	fmt.Println()
	fmt.Println(c(cyan+bold, divider))
	fmt.Println(c(cyan+bold, "  "+icon+"  "+strings.ToUpper(title)))
	fmt.Println(c(cyan+bold, divider))
}

// withSpinner runs fn while animating a spinner on the current line, then
// clears the line. The caller is responsible for printing the result line.
func withSpinner(label string, fn func() error) (time.Duration, error) {
	start := time.Now()
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		i := 0
		t := time.NewTicker(90 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-stop:
				return
			case <-t.C:
				frame := spinnerFrames[i%len(spinnerFrames)]
				fmt.Printf("\r%s %s", c(cyan, frame), label)
				i++
			}
		}
	}()
	err := fn()
	close(stop)
	<-done
	fmt.Print("\r\033[K")
	return time.Since(start), err
}

// ─────────────────────────────────────────────────────────────────────────
// Shell execution helpers
// ─────────────────────────────────────────────────────────────────────────

var currentLog *os.File

func shRaw(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	cmd.Stdin = nil
	err := cmd.Run()
	return buf.String(), err
}

// sh runs a shell command, logs it if a per-tool log file is open, and
// transparently retries pip installs that hit Debian's PEP-668
// "externally-managed-environment" guard with --break-system-packages.
func sh(cmdStr string) (string, error) {
	out, err := shRaw(cmdStr)
	if currentLog != nil {
		fmt.Fprintf(currentLog, "$ %s\n%s\n", cmdStr, out)
	}
	if err != nil && strings.Contains(out, "externally-managed-environment") &&
		(strings.Contains(cmdStr, "pip3 install") || strings.Contains(cmdStr, "pip install")) &&
		!strings.Contains(cmdStr, "--break-system-packages") {
		retry := strings.Replace(cmdStr, "install", "install --break-system-packages", 1)
		rOut, rErr := shRaw(retry)
		if currentLog != nil {
			fmt.Fprintf(currentLog, "$ %s\n%s\n", retry, rOut)
		}
		return rOut, rErr
	}
	return out, err
}

func errFrom(out string, err error) error {
	if err == nil {
		return nil
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	last := ""
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			last = strings.TrimSpace(lines[i])
			break
		}
	}
	if last == "" {
		return err
	}
	if len(last) > 160 {
		last = last[:160] + "…"
	}
	return fmt.Errorf("%s", last)
}

func shortErr(err error) string {
	s := err.Error()
	if len(s) > 100 {
		s = s[:100] + "…"
	}
	return s
}

func sanitize(name string) string {
	r := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ".", "_")
	return r.Replace(strings.ToLower(name))
}

var sudoCmd string    // "" if already root, otherwise "sudo "
var toolsDir string   // where git-cloned source tools live
var userBinDir string // where wrapper scripts / user installs live
var logsDir string

// ─────────────────────────────────────────────────────────────────────────
// Install-method helpers — each returns a func() error usable as Tool.Install
// ─────────────────────────────────────────────────────────────────────────

func goInstall(pkg string) func() error {
	return func() error {
		out, err := sh("go install -v " + pkg + "@latest")
		return errFrom(out, err)
	}
}

func goInstallVersion(pkg, ver string) func() error {
	return func() error {
		out, err := sh("go install -v " + pkg + "@" + ver)
		return errFrom(out, err)
	}
}

func apt(pkgs ...string) func() error {
	return func() error {
		out, err := sh(sudoCmd + "apt-get install -y " + strings.Join(pkgs, " "))
		return errFrom(out, err)
	}
}

func aptBestEffort(pkgs ...string) {
	sh(sudoCmd + "apt-get install -y " + strings.Join(pkgs, " "))
}

func pipUser(pkgs ...string) func() error {
	return func() error {
		out, err := sh("pip3 install --user " + strings.Join(pkgs, " "))
		return errFrom(out, err)
	}
}

func pipxTool(pkg string) func() error {
	return func() error {
		out, err := sh("pipx install " + pkg)
		return errFrom(out, err)
	}
}

func cargoTool(pkg string) func() error {
	return func() error {
		out, err := sh("cargo install " + pkg)
		return errFrom(out, err)
	}
}

func npmGlobal(pkg string) func() error {
	return func() error {
		out, err := sh(sudoCmd + "npm install -g " + pkg)
		return errFrom(out, err)
	}
}

func rawCmd(cmdline string) func() error {
	return func() error {
		out, err := sh(cmdline)
		return errFrom(out, err)
	}
}

// gitMake clones repo (if not already present) into toolsDir/dirName, runs
// buildCmd inside it, then locates the resulting binary by name anywhere
// under the checkout and symlinks it into userBinDir.
func gitMake(repo, dirName, buildCmd, binName string) func() error {
	return func() error {
		dest := filepath.Join(toolsDir, dirName)
		if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
			out, err := sh("git clone --depth=1 " + repo + " " + dest)
			if err != nil {
				return errFrom(out, err)
			}
		}
		out, err := sh("cd " + dest + " && " + buildCmd)
		if err != nil {
			return errFrom(out, err)
		}
		findOut, _ := sh("find " + dest + " -maxdepth 4 -type f -name " + binName + " -perm -u+x 2>/dev/null | head -n1")
		found := strings.TrimSpace(findOut)
		if found == "" {
			return fmt.Errorf("build finished but binary '%s' was not found under %s", binName, dest)
		}
		linkDest := filepath.Join(userBinDir, binName)
		os.Remove(linkDest)
		return os.Symlink(found, linkDest)
	}
}

// gitPy clones a Python tool and wraps its entry script with a small shim on
// PATH so it behaves like a normal installed command.
func gitPy(repo, dirName, entryScript, setupCmd string) func() error {
	return func() error {
		dest := filepath.Join(toolsDir, dirName)
		if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
			out, err := sh("git clone --depth=1 " + repo + " " + dest)
			if err != nil {
				return errFrom(out, err)
			}
		}
		if setupCmd != "" {
			out, err := sh("cd " + dest + " && " + setupCmd)
			if err != nil {
				return errFrom(out, err)
			}
		}
		wrapperName := strings.ToLower(dirName)
		wrapper := "#!/usr/bin/env bash\nexec python3 \"" + filepath.Join(dest, entryScript) + "\" \"$@\"\n"
		if err := os.WriteFile(filepath.Join(userBinDir, wrapperName), []byte(wrapper), 0o755); err != nil {
			return err
		}
		return nil
	}
}

func anyBin(names ...string) func() bool {
	return func() bool {
		for _, n := range names {
			if _, err := exec.LookPath(n); err == nil {
				return true
			}
		}
		return false
	}
}

func pyModuleInstalled(module string) func() bool {
	return func() bool {
		_, err := sh("python3 -c 'import " + module + "'")
		return err == nil
	}
}

func dpkgInstalled(pkg string) func() bool {
	return func() bool {
		_, err := sh("dpkg -s " + pkg)
		return err == nil
	}
}
