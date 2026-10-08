<div align="center">

# bb-arsenal

**One command to set up a bug bounty / recon workstation.**

Installs ~100 recon, fuzzing, scanning, and analysis tools, skips whatever you already have, and shows progress in a clean terminal UI.

[![CI](https://github.com/BlackRootSec/bb-arsenal/actions/workflows/ci.yml/badge.svg)](https://github.com/BlackRootSec/bb-arsenal/actions/workflows/ci.yml)
[![Go version](https://img.shields.io/github/go-mod/go-version/BlackRootSec/bb-arsenal)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## Features

- **99 tools + 13 prerequisites** across 19 categories: subdomain enumeration, DNS, HTTP/port scanning, crawling, fuzzing, parameter discovery, vulnerability scanning, XSS, SQLi, API/JWT, secrets, cloud, mobile/RE, OSINT, and more.
- **Skips what's already installed.** Safe to re-run any time.
- **Right installer per tool:** `go install`, `apt`, `pipx`/`pip`, `cargo`, `npm`, vendor scripts, or `git clone` + build.
- **Handles Debian's locked-down Python** (PEP 668) automatically, so pip-based tools don't die with `externally-managed-environment`.
- **Nothing is a dead end:** every install's full output goes to `~/.bb-installer/logs/<tool>.log`; the terminal shows a one-line reason.
- **Zero dependencies:** a single Go binary built from the standard library only.
- **Filterable:** run everything, a few categories, or just specific tools.

## Install

**With Go (1.21+):**

```bash
go install github.com/BlackRootSec/bb-arsenal@latest
```

**From source:**

```bash
git clone https://github.com/BlackRootSec/bb-arsenal.git
cd bb-arsenal
make build            # or: go build -o bb-arsenal .
```

**Prebuilt binary:** grab the Linux `amd64` / `arm64` archive from the [Releases](https://github.com/BlackRootSec/bb-arsenal/releases) page.

## Usage

Run it as your **normal user** (it calls `sudo` itself for `apt`/`snap`), or as root on a dedicated box. Don't prefix it with `sudo`: user-level tools (Go, pipx, cargo) would land in root's home instead of yours, and the installer will warn you if it detects this.

```bash
bb-arsenal                      # install everything (asks for confirmation first)
bb-arsenal -y                   # ...without the prompt
bb-arsenal --dry-run            # show what would be installed, change nothing
bb-arsenal --list               # print the full catalog and exit

bb-arsenal --categories subs,dns,vuln       # only these categories
bb-arsenal --only nuclei,httpx,subfinder    # only these tools
bb-arsenal --skip masscan,rustscan          # everything except these
bb-arsenal --no-gui                         # skip Burp/Wireshark/Postman/Insomnia
bb-arsenal --no-color                       # plain output (CI logs, pipes)
bb-arsenal --version
```

| Flag | Description |
|------|-------------|
| `-y`, `--yes` | Don't ask for confirmation |
| `--dry-run` | Report what would happen; install nothing |
| `--list` | Show the catalog (respects the filters below) and exit |
| `--categories a,b` | Run only these category keys (table below) |
| `--only a,b` | Run only these tools (exact names from `--list`, case-insensitive) |
| `--skip a,b` | Exclude these tools |
| `--no-gui` | Skip the optional GUI-apps category |
| `--no-color` | Disable ANSI colors |
| `--version` | Print version and exit |

> **Tip:** the *Prerequisites* category (git, Go, pipx, cargo, npm...) is a normal category. On a fresh machine, run once without filters (or with `--categories prereq`) before using `--only`.

### What the output looks like

Abridged example:

```
──────────────────────────────────────────────────────────────
  🌐  SUBDOMAIN ENUMERATION
──────────────────────────────────────────────────────────────
[ 14/112] ✔ subfinder                        installed in 21.4s — Passive sources need API keys
[ 15/112] ✔ amass                            already installed — skipped
[ 16/112] ⠹ installing assetfinder …
[ 17/112] ✘ findomain                        failed: error: could not compile ...  (log: ~/.bb-installer/logs/findomain.log)
```

## What gets installed

| Key | Category | Tools |
|-----|----------|-------|
| `prereq` | Prerequisites | `git`, `curl`, `wget`, `unzip`, `build-essential (gcc/make)`, `python3`, `python3-pip`, `pipx`, `Go toolchain`, `Rust/Cargo`, `Node.js/npm`, `ripgrep`, `jq` |
| `subs` | Subdomain Enumeration | `subfinder`, `amass`, `assetfinder`, `findomain`, `chaos`, `github-subdomains` |
| `dns` | DNS Resolution & Permutation | `dnsx`, `massdns`, `puredns`, `shuffledns`, `alterx`, `dnsgen`, `altdns` |
| `http` | HTTP Probing & Port Scanning | `httpx`, `naabu`, `nmap`, `uncover`, `masscan`, `rustscan` |
| `crawl` | Crawling & URL Discovery | `katana`, `gau`, `waybackurls`, `gospider`, `hakrawler`, `urlfinder`, `urlhunter`, `cariddi`, `xnLinkFinder`, `subjs`, `LinkFinder`, `SecretFinder`, `JSluice` |
| `fuzz` | Fuzzing & Content Discovery | `ffuf`, `feroxbuster`, `gobuster`, `dirsearch`, `wfuzz`, `dirb`, `kiterunner` |
| `params` | Parameter Discovery | `arjun`, `x8`, `ParamSpider`, `unfurl`, `qsreplace`, `uro` |
| `vuln` | Vulnerability Scanning | `nuclei`, `nuclei-templates`, `nikto`, `testssl.sh` |
| `xss` | XSS Testing | `dalfox`, `kxss`, `Gxss`, `XSStrike` |
| `sqli` | SQL Injection | `sqlmap`, `ghauri` |
| `api` | API, Proxy & Protocol Tools | `httpie`, `mitmproxy`, `interactsh-client`, `InQL` |
| `jwt` | JWT Tooling | `jwt_tool`, `PyJWT` |
| `secrets` | Secrets & Static Analysis | `trufflehog`, `gitleaks`, `git-secrets`, `detect-secrets`, `semgrep`, `Retire.js` |
| `tech` | Tech Fingerprinting | `whatweb`, `webanalyze` |
| `cloud` | Cloud Security | `cloud_enum`, `CloudBrute`, `Prowler`, `ScoutSuite`, `Pacu`, `S3Scanner` |
| `mobile` | Mobile & Reverse Engineering | `apktool`, `adb`, `jadx`, `frida-tools`, `objection`, `radare2`, `apk-mitm`, `binutils (strings/objdump/readelf/nm)` |
| `netutils` | Network Utilities | `socat`, `netcat-openbsd`, `tcpdump`, `tshark`, `dnsutils`, `whois`, `net-tools`, `proxychains4` |
| `osint` | OSINT | `theHarvester`, `SpiderFoot`, `Shodan CLI`, `Censys CLI` |
| `gui` | Optional GUI Apps | `Burp Suite Community`, `Wireshark`, `Postman`, `Insomnia` |

Install methods across all 112 entries (tools + prerequisites): apt 36, go 35, pipx 15, git 13, cargo 4, pip 4, npm 2, snap 2, vendor script 1.

## Where things go

| Path | Contents |
|------|----------|
| `~/go/bin` | `go install` binaries |
| `~/.cargo/bin` | `cargo install` binaries |
| `~/.local/bin` | pipx apps, plus small launcher shims for git-cloned Python tools and symlinks for git-built binaries |
| `~/bbtools/` | git checkouts (massdns, LinkFinder, XSStrike, jwt_tool, ...) |
| `~/nuclei-templates` | nuclei template repo |
| `~/.bb-installer/logs/` | one log per tool install attempt |

To make everything reachable, the installer adds those bin directories to `PATH` for its own run and appends **one** line (marked `# added by bb-installer`) to `~/.bashrc` and `~/.zshrc` **if those files already exist**. It's idempotent and won't append twice. Run `source ~/.bashrc` or open a new terminal afterwards.

## Not included, and why

The wishlist this project grew from mixed real tools with things that aren't software. `bb-arsenal --list` prints the same explanation; in short:

- **Protocols and specs** (REST, GraphQL, WebSockets, gRPC, SOAP, OpenAPI, Swagger) and **vulnerability classes** (OAuth/OIDC, JWT, API keys, BOLA/IDOR, BFLA, mass assignment, rate limiting, race conditions) are things you *test for*, not install. The HTTP/API/JWT tools above are what you test them with.
- **Websites and SaaS** (crt.sh, jwt.io, Shodan, Censys, FOFA, ZoomEye, SecurityTrails): the Shodan and Censys CLIs are installed; you still need your own accounts and API keys.
- **Burp BApp Store extensions** (DOM Invader, Autorize, AuthMatrix, JWT Editor, Param Miner, HTTP Request Smuggler) install from inside Burp.
- **Heavyweight / GUI / Windows-only** (Ghidra, Cutter, MobSF, x64dbg, WinDbg) need a JDK, Docker, or Windows; the installer prints the right pointer instead of a fragile automated install.

## Platform support

Built for **Debian-family Linux** (Kali, Parrot, Ubuntu, Debian), because roughly a third of the catalog installs through `apt`. On other distros the `go`/`cargo`/`pipx`/`npm` tools can still work, but `apt`-based entries (and the prerequisites step) will fail; use `--skip` or install those natively. The GUI category uses `snap` for Postman and Insomnia, so use `--no-gui` on systems without snap.

## Status and accuracy

Install commands (including versioned Go module paths such as `/v2`, `/v3`) were checked against each project's current upstream docs, but upstream projects rename binaries and change install methods over time, and not every tool has been exercised on every distro. If something breaks, check the log first, then [open an issue](https://github.com/BlackRootSec/bb-arsenal/issues/new/choose) with the log excerpt. Fixes are usually a one-line edit; see [CONTRIBUTING.md](CONTRIBUTING.md).

## Troubleshooting

- **A tool failed:** read `~/.bb-installer/logs/<tool>.log`. The terminal line only shows the last error line.
- **`go install` says a newer Go is required:** with Go 1.21+ and internet access the toolchain is fetched automatically; otherwise install a current Go from <https://go.dev/dl/>.
- **`cargo install` fails to compile:** the distro's Rust can be too old for some crates. Install via [rustup](https://rustup.rs), then re-run (already-installed tools are skipped).
- **Installed tool "not found":** `source ~/.bashrc` or open a new terminal.
- **kiterunner build fails:** use a prebuilt binary from <https://github.com/assetnote/kiterunner/releases>.

## Responsible use

These are offensive-security tools. Only run them against systems you own or have **explicit written authorization** to test (for example, assets in scope of a bug bounty program, following that program's rules). You are responsible for how you use them. `bb-arsenal` only *installs* third-party software; each tool is governed by its own license and terms.

## Contributing

Adding or fixing a tool is a small, local change. See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE) © 2026 BlackRootSec
