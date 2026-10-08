package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	flagList       bool
	flagYes        bool
	flagDryRun     bool
	flagNoColor    bool
	flagNoGUI      bool
	flagVersion    bool
	flagCategories string
	flagOnly       string
	flagSkip       string
)

func parseFlags() {
	flag.BoolVar(&flagList, "list", false, "list the full tool catalog and exit without installing anything")
	flag.BoolVar(&flagYes, "y", false, "don't ask for confirmation before installing")
	flag.BoolVar(&flagYes, "yes", false, "don't ask for confirmation before installing")
	flag.BoolVar(&flagDryRun, "dry-run", false, "show what would be installed without doing it")
	flag.BoolVar(&flagVersion, "version", false, "print version and exit")
	flag.BoolVar(&flagNoColor, "no-color", false, "disable ANSI colors")
	flag.BoolVar(&flagNoGUI, "no-gui", false, "skip the optional GUI-apps category (Burp/Wireshark/Postman/Insomnia)")
	flag.StringVar(&flagCategories, "categories", "", "comma-separated category keys to run (default: all). See --list for keys.")
	flag.StringVar(&flagOnly, "only", "", "comma-separated tool names to install, ignoring everything else (exact names as shown by --list, case-insensitive)")
	flag.StringVar(&flagSkip, "skip", "", "comma-separated tool names to exclude (exact names as shown by --list, case-insensitive)")
	flag.Parse()
	noColor = flagNoColor
}

func csvSet(s string) map[string]bool {
	m := map[string]bool{}
	for _, p := range strings.Split(s, ",") {
		p = strings.ToLower(strings.TrimSpace(p))
		if p != "" {
			m[p] = true
		}
	}
	return m
}

func setupDirs() {
	home := homeDir()
	toolsDir = filepath.Join(home, "bbtools")
	userBinDir = filepath.Join(home, ".local", "bin")
	logsDir = filepath.Join(home, ".bb-installer", "logs")
}

// makeDirs creates the working directories. Only called for real installs so
// that --list and --dry-run leave the filesystem untouched.
func makeDirs() {
	os.MkdirAll(toolsDir, 0o755)
	os.MkdirAll(userBinDir, 0o755)
	os.MkdirAll(logsDir, 0o755)
}

func setupPath() []string {
	home := homeDir()
	extras := []string{
		filepath.Join(home, "go", "bin"),
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, ".cargo", "bin"),
		"/usr/local/bin",
		"/snap/bin",
	}
	cur := os.Getenv("PATH")
	existing := map[string]bool{}
	for _, p := range strings.Split(cur, ":") {
		existing[p] = true
	}
	for _, e := range extras {
		if !existing[e] {
			cur += ":" + e
		}
	}
	os.Setenv("PATH", cur)
	return extras
}

func ensureRCHasPath(dirs []string) {
	home := homeDir()
	marker := "# added by bb-installer"
	line := "export PATH=\"$PATH:" + strings.Join(dirs, ":") + "\"  " + marker
	for _, rc := range []string{filepath.Join(home, ".bashrc"), filepath.Join(home, ".zshrc")} {
		data, err := os.ReadFile(rc)
		if err != nil {
			continue
		}
		if strings.Contains(string(data), marker) {
			continue
		}
		f, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			continue
		}
		fmt.Fprintf(f, "\n%s\n", line)
		f.Close()
	}
}

func sudoWarmup() {
	if sudoCmd == "" {
		return
	}
	fmt.Println(c(yellow, "This installer needs sudo for system packages (apt/snap). You may be prompted for your password."))
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	go func() {
		for {
			time.Sleep(60 * time.Second)
			exec.Command("sudo", "-n", "true").Run()
		}
	}()
}

type counters struct {
	installed, skipped, failed, manual int
	failedNames                        []string
}

func processTool(t Tool, idx, total int, cnt *counters) {
	prefix := fmt.Sprintf("[%3d/%3d]", idx, total)
	if t.Manual != "" {
		fmt.Printf("%s %s %-32s %s\n", c(gray, prefix), c(yellow, "○"), t.Name, c(gray, t.Manual))
		cnt.manual++
		return
	}

	installed := false
	if t.CheckFunc != nil {
		installed = t.CheckFunc()
	} else if t.Bin != "" {
		_, err := exec.LookPath(t.Bin)
		installed = err == nil
	}

	if installed {
		fmt.Printf("%s %s %-32s %s\n", c(gray, prefix), c(green, "✔"), t.Name, c(gray, "already installed — skipped"))
		cnt.skipped++
		return
	}

	if flagDryRun {
		fmt.Printf("%s %s %-32s %s\n", c(gray, prefix), c(cyan, "→"), t.Name, c(gray, "would install ("+t.Method+")"))
		return
	}

	if t.Install == nil {
		fmt.Printf("%s %s %-32s %s\n", c(gray, prefix), c(yellow, "○"), t.Name, c(gray, "no installer defined — skipped"))
		cnt.manual++
		return
	}

	label := fmt.Sprintf("%s %s installing %s …", c(gray, prefix), c(cyan, "◌"), t.Name)
	logPath := filepath.Join(logsDir, sanitize(t.Name)+".log")
	f, _ := os.Create(logPath)
	currentLog = f
	dur, err := withSpinner(label, t.Install)
	if f != nil {
		f.Close()
	}
	currentLog = nil

	if err != nil {
		fmt.Printf("%s %s %-32s %s\n", c(gray, prefix), c(red, "✘"), t.Name,
			c(red, "failed: "+shortErr(err))+c(gray, "  (log: "+logPath+")"))
		cnt.failed++
		cnt.failedNames = append(cnt.failedNames, t.Name)
		return
	}

	extra := ""
	if t.Info != "" {
		extra = c(gray, " — "+t.Info)
	}
	fmt.Printf("%s %s %-32s %s%s\n", c(gray, prefix), c(green, "✔"), t.Name,
		c(green, fmt.Sprintf("installed in %s", dur.Round(100*time.Millisecond))), extra)
	cnt.installed++
}

func printManualSection(title string, notes []ManualNote) {
	if len(notes) == 0 {
		return
	}
	fmt.Println()
	fmt.Println(c(yellow+bold, "  "+title))
	for _, n := range notes {
		fmt.Println(c(yellow, "   • "+n.Name))
		fmt.Println(c(gray, "     "+n.Reason))
	}
}

func printCatalog(categories []Category) {
	for _, cat := range categories {
		sectionHeader(cat.Icon, cat.Title)
		for _, t := range cat.Tools {
			method := t.Method
			if method == "" {
				method = "manual"
			}
			fmt.Printf("  %-38s %s\n", t.Name, c(gray, "["+method+"]"))
		}
	}
	fmt.Println()
	fmt.Println(c(bold, "Items intentionally not included (concepts, web services, Burp extensions, etc.):"))
	printManualSection("GUI / heavyweight tools", manualReconTools())
	printManualSection("Burp Suite BApp Store extensions", manualBurpExtensions())
	printManualSection("Websites & SaaS platforms", manualWebServices())
	printManualSection("Protocols & vulnerability classes (not software)", manualConcepts())
}

func main() {
	parseFlags()
	if flagVersion {
		fmt.Printf("bb-arsenal %s (%s)\n", version, repoURL)
		return
	}
	setupDirs()
	pathDirs := setupPath()

	if os.Geteuid() == 0 {
		sudoCmd = ""
	} else {
		sudoCmd = "sudo "
	}

	if !flagList && os.Geteuid() == 0 && os.Getenv("SUDO_USER") != "" {
		fmt.Println(c(yellow, "⚠ Running via sudo: user-level tools (go/pipx/cargo/git) will be installed into root's home,"))
		fmt.Println(c(yellow, "  not "+os.Getenv("SUDO_USER")+"'s. Prefer running as your normal user — the installer calls sudo itself where needed."))
	}

	categories := buildCategories()
	if flagNoGUI {
		filtered := categories[:0]
		for _, cat := range categories {
			if cat.Key != "gui" {
				filtered = append(filtered, cat)
			}
		}
		categories = filtered
	}

	catFilter := csvSet(flagCategories)
	onlyFilter := csvSet(flagOnly)
	skipFilter := csvSet(flagSkip)

	if len(catFilter) > 0 {
		filtered := categories[:0]
		for _, cat := range categories {
			if catFilter[cat.Key] {
				filtered = append(filtered, cat)
			}
		}
		categories = filtered
	}

	// Apply --only / --skip at the tool level.
	var finalCats []Category
	for _, cat := range categories {
		var kept []Tool
		for _, t := range cat.Tools {
			name := strings.ToLower(t.Name)
			if len(onlyFilter) > 0 && !onlyFilter[name] {
				continue
			}
			if skipFilter[name] {
				continue
			}
			kept = append(kept, t)
		}
		if len(kept) > 0 {
			cat.Tools = kept
			finalCats = append(finalCats, cat)
		}
	}
	categories = finalCats

	if len(categories) == 0 {
		fmt.Println(c(yellow, "No tools matched your --categories/--only/--skip filters. Try --list to see valid keys and names."))
		return
	}

	if flagList {
		banner()
		printCatalog(categories)
		return
	}

	total := 0
	for _, cat := range categories {
		total += len(cat.Tools)
	}

	banner()
	if runtime.GOOS != "linux" {
		fmt.Println(c(yellow, "⚠ This installer targets Debian/Ubuntu/Kali Linux. apt-based tools will likely fail on "+runtime.GOOS+"."))
	}
	fmt.Printf("%s\n", c(gray, fmt.Sprintf("Tools directory: %s   |   User bin: %s   |   Logs: %s", toolsDir, userBinDir, logsDir)))
	fmt.Printf("%s %d tools across %d categories.\n", c(bold, "Ready to process"), total, len(categories))

	if !flagYes && !flagDryRun {
		fmt.Print(c(yellow, "\nContinue? [Y/n] "))
		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		resp = strings.ToLower(strings.TrimSpace(resp))
		if resp == "n" || resp == "no" {
			fmt.Println("Aborted — nothing was installed.")
			return
		}
	}

	if !flagDryRun {
		makeDirs()
		ensureRCHasPath(pathDirs)
		sudoWarmup()
	}

	cnt := &counters{}
	idx := 0
	for _, cat := range categories {
		sectionHeader(cat.Icon, cat.Title)
		for _, t := range cat.Tools {
			idx++
			processTool(t, idx, total, cnt)
		}
	}

	fmt.Println()
	fmt.Println(c(cyan+bold, divider))
	fmt.Println(c(cyan+bold, "  SUMMARY"))
	fmt.Println(c(cyan+bold, divider))
	fmt.Printf("  %s Installed:       %d\n", c(green, "✔"), cnt.installed)
	fmt.Printf("  %s Already present: %d\n", c(gray, "●"), cnt.skipped)
	fmt.Printf("  %s Failed:          %d\n", c(red, "✘"), cnt.failed)
	fmt.Printf("  %s Manual steps:    %d\n", c(yellow, "○"), cnt.manual)

	if len(cnt.failedNames) > 0 {
		fmt.Println()
		fmt.Println(c(red, "  Failed — check the per-tool logs under "+logsDir+":"))
		for _, n := range cnt.failedNames {
			fmt.Println(c(red, "    • "+n))
		}
	}

	fmt.Println()
	fmt.Println(c(bold, "A few items from a typical wishlist aren't \"installable tools\" — here's why, and what covers them instead:"))
	printManualSection("GUI / heavyweight tools", manualReconTools())
	printManualSection("Burp Suite BApp Store extensions", manualBurpExtensions())
	printManualSection("Websites & SaaS platforms", manualWebServices())
	printManualSection("Protocols & vulnerability classes (not software)", manualConcepts())

	if !flagDryRun {
		fmt.Println()
		fmt.Println(c(gray, "PATH entries (go/bin, .local/bin, .cargo/bin) were added to ~/.bashrc and ~/.zshrc if those files exist."))
		fmt.Println(c(gray, "Run `source ~/.bashrc` (or open a new terminal) so every installed tool is on your PATH."))
	}
}
