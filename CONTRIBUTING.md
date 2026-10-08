# Contributing to bb-arsenal

Thanks for helping keep the catalog accurate! Most contributions are a one-line fix or a new `Tool{...}` entry.

## Project layout

| File | Purpose |
|------|---------|
| `types.go` | `Tool` / `Category` types |
| `catalog.go` | `buildCategories()`: category order, keys, titles, icons |
| `catalog_recon.go` | Prerequisites, subdomains, DNS, HTTP/ports, crawling, fuzzing, params, vuln scanning |
| `catalog_testing.go` | XSS, SQLi, API, JWT, secrets, tech, cloud, mobile, network, OSINT, GUI |
| `notes.go` | "Not installable" explanations shown by `--list` and at the end of a run |
| `main.go` | Styling, spinner, shell helpers, install-method helpers |
| `run.go` | Flags, PATH/sudo setup, the install loop, summary |
| `catalog_test.go` | Catalog integrity tests |

## Adding a tool

1. Pick the matching category function (or propose a new category in `catalog.go`).
2. Add an entry:

   ```go
   {Name: "subfinder", Method: "go", Bin: "subfinder",
       Install: goInstall("github.com/projectdiscovery/subfinder/v2/cmd/subfinder"),
       Info:    "Passive sources need API keys"},
   ```

   | Field | Meaning |
   |-------|---------|
   | `Name` | Display name. Must be unique across the whole catalog. `--only`/`--skip` match it case-insensitively. |
   | `Method` | Label shown by `--list`: `go`, `apt`, `pip`, `pipx`, `cargo`, `npm`, `git`, `script`, `snap`. |
   | `Bin` | Binary looked up on `PATH` to decide "already installed". Must be the name **after** install. |
   | `CheckFunc` | Use instead of `Bin` when there's no binary (a Python library, a templates repo, ...). Helpers: `anyBin`, `pyModuleInstalled`, `dpkgInstalled`. |
   | `Install` | A `func() error`. Prefer a helper (below). |
   | `Info` | Optional one-liner printed after success (API keys to configure, etc.). |

3. Install helpers in `main.go`:

   | Helper | Runs |
   |--------|------|
   | `goInstall(pkg)` / `goInstallVersion(pkg, ver)` | `go install -v pkg@latest` / `@ver` |
   | `apt(pkgs...)` | `apt-get install -y ...` (adds `sudo` when not root) |
   | `pipxTool(pkg)` | `pipx install pkg`, preferred for Python CLIs |
   | `pipUser(pkgs...)` | `pip3 install --user ...` (auto-retries with `--break-system-packages` on PEP 668 systems) |
   | `cargoTool(crate)` | `cargo install crate` |
   | `npmGlobal(pkg)` | `npm install -g pkg` |
   | `gitMake(repo, dir, buildCmd, bin)` | clone, build, symlink the built binary into `~/.local/bin` |
   | `gitPy(repo, dir, entry.py, setupCmd)` | clone a Python tool and create a launcher shim on `PATH` |
   | `rawCmd(cmdline)` | any shell command (vendor install scripts, snap, ...) |

4. Use the **upstream project's own** documented install command. Watch for versioned Go module paths (`/v2`, `/v3`); a wrong suffix is the most common breakage.

5. Verify in a clean container (e.g. `docker run -it kalilinux/kali-rolling`): confirm the install works and that `Bin` is what actually lands on `PATH`.

## Before you open a PR

```bash
make check     # gofmt, go vet, go test (same as CI)
make build && ./bb-arsenal --list
```

`catalog_test.go` fails if a tool has no installer, no skip-check, a missing `Method`, or a duplicate name. Fix those before pushing.

## Commit style

Short, imperative subject lines, e.g. `fix: correct naabu module path`, `feat: add trufflehog`, `docs: clarify sudo usage`. The `docs:`, `test:`, and `chore:` prefixes are omitted from release notes.

## Scope

In scope: command-line tools used in authorized recon / bug bounty / appsec workflows that install cleanly via a standard mechanism. Out of scope: protocols, vulnerability classes, websites/SaaS, and anything that can't be installed non-interactively; those belong in `notes.go`.
