package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ── Prerequisites ──────────────────────────────────────────────────────────

func prereqTools() []Tool {
	return []Tool{
		{Name: "git", Method: "apt", Bin: "git", Install: apt("git")},
		{Name: "curl", Method: "apt", Bin: "curl", Install: apt("curl")},
		{Name: "wget", Method: "apt", Bin: "wget", Install: apt("wget")},
		{Name: "unzip", Method: "apt", Bin: "unzip", Install: apt("unzip")},
		{Name: "build-essential (gcc/make)", Method: "apt", Bin: "make", Install: apt("build-essential")},
		{Name: "python3", Method: "apt", Bin: "python3", Install: apt("python3")},
		{Name: "python3-pip", Method: "apt", Bin: "pip3", Install: apt("python3-pip")},
		{Name: "pipx", Method: "apt", Bin: "pipx", Install: func() error {
			out, err := sh(sudoCmd + "apt-get install -y pipx")
			if err == nil {
				sh("pipx ensurepath")
			}
			return errFrom(out, err)
		}},
		{Name: "Go toolchain", Method: "apt", Bin: "go", Install: apt("golang-go"),
			Info: "If a tool needs a newer Go than apt provides, `go install` will auto-fetch the right toolchain version as long as you have internet access."},
		{Name: "Rust/Cargo", Method: "apt", Bin: "cargo", Install: apt("cargo")},
		{Name: "Node.js/npm", Method: "apt", Bin: "npm", Install: apt("nodejs", "npm")},
		{Name: "ripgrep", Method: "apt", Bin: "rg", Install: apt("ripgrep")},
		{Name: "jq", Method: "apt", Bin: "jq", Install: apt("jq")},
	}
}

// ── Subdomain enumeration ──────────────────────────────────────────────────

func subdomainTools() []Tool {
	return []Tool{
		{Name: "subfinder", Method: "go", Bin: "subfinder",
			Install: goInstall("github.com/projectdiscovery/subfinder/v2/cmd/subfinder"),
			Info:    "Passive sources need API keys — see ~/.config/subfinder/provider-config.yaml"},
		{Name: "amass", Method: "go", Bin: "amass",
			Install: goInstallVersion("github.com/owasp-amass/amass/v4/...", "master")},
		{Name: "assetfinder", Method: "go", Bin: "assetfinder",
			Install: goInstall("github.com/tomnomnom/assetfinder")},
		{Name: "findomain", Method: "cargo", Bin: "findomain",
			Install: cargoTool("findomain")},
		{Name: "chaos", Method: "go", Bin: "chaos",
			Install: goInstall("github.com/projectdiscovery/chaos-client/cmd/chaos"),
			Info:    "Needs a PDCP/Chaos API key: https://chaos.projectdiscovery.io"},
		{Name: "github-subdomains", Method: "go", Bin: "github-subdomains",
			Install: goInstall("github.com/gwen001/github-subdomains"),
			Info:    "Needs a GitHub token in the GITHUB_TOKEN(S) env var"},
	}
}

// ── DNS resolution, brute force & permutation ─────────────────────────────

func dnsTools() []Tool {
	return []Tool{
		{Name: "dnsx", Method: "go", Bin: "dnsx",
			Install: goInstall("github.com/projectdiscovery/dnsx/cmd/dnsx")},
		{Name: "massdns", Method: "git", Bin: "massdns",
			Install: gitMake("https://github.com/blechschmidt/massdns.git", "massdns", "make", "massdns")},
		{Name: "puredns", Method: "go", Bin: "puredns",
			Install: goInstall("github.com/d3mondev/puredns/v2"),
			Info:    "Needs massdns on PATH and a good resolvers list to work well"},
		{Name: "shuffledns", Method: "go", Bin: "shuffledns",
			Install: goInstall("github.com/projectdiscovery/shuffledns/cmd/shuffledns")},
		{Name: "alterx", Method: "go", Bin: "alterx",
			Install: goInstall("github.com/projectdiscovery/alterx/cmd/alterx")},
		{Name: "dnsgen", Method: "pip", Bin: "dnsgen",
			Install: pipUser("dnsgen")},
		{Name: "altdns", Method: "pip", Bin: "altdns",
			Install: pipUser("py-altdns")},
	}
}

// ── HTTP probing & port scanning ──────────────────────────────────────────

func httpPortTools() []Tool {
	return []Tool{
		{Name: "httpx", Method: "go", Bin: "httpx",
			Install: goInstall("github.com/projectdiscovery/httpx/cmd/httpx")},
		{Name: "naabu", Method: "go", Bin: "naabu",
			Install: func() error {
				aptBestEffort("libpcap-dev")
				out, err := sh("go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest")
				return errFrom(out, err)
			}},
		{Name: "nmap", Method: "apt", Bin: "nmap", Install: apt("nmap")},
		{Name: "uncover", Method: "go", Bin: "uncover",
			Install: goInstall("github.com/projectdiscovery/uncover/cmd/uncover"),
			Info:    "Needs API keys for Shodan/Censys/Fofa/etc in ~/.config/uncover/provider-config.yaml"},
		{Name: "masscan", Method: "apt", Bin: "masscan", Install: apt("masscan")},
		{Name: "rustscan", Method: "cargo", Bin: "rustscan", Install: cargoTool("rustscan")},
	}
}

// ── Crawling, JS analysis & URL discovery ─────────────────────────────────

func crawlTools() []Tool {
	return []Tool{
		{Name: "katana", Method: "go", Bin: "katana",
			Install: goInstall("github.com/projectdiscovery/katana/cmd/katana")},
		{Name: "gau", Method: "go", Bin: "gau",
			Install: goInstall("github.com/lc/gau/v2/cmd/gau")},
		{Name: "waybackurls", Method: "go", Bin: "waybackurls",
			Install: goInstall("github.com/tomnomnom/waybackurls")},
		{Name: "gospider", Method: "go", Bin: "gospider",
			Install: goInstall("github.com/jaeles-project/gospider")},
		{Name: "hakrawler", Method: "go", Bin: "hakrawler",
			Install: goInstall("github.com/hakluke/hakrawler")},
		{Name: "urlfinder", Method: "go", Bin: "urlfinder",
			Install: goInstall("github.com/projectdiscovery/urlfinder/cmd/urlfinder")},
		{Name: "urlhunter", Method: "go", Bin: "urlhunter",
			Install: goInstall("github.com/utkusen/urlhunter")},
		{Name: "cariddi", Method: "go", Bin: "cariddi",
			Install: goInstall("github.com/edoardottt/cariddi/cmd/cariddi")},
		{Name: "xnLinkFinder", Method: "pipx", Bin: "xnLinkFinder",
			Install: pipxTool("xnLinkFinder")},
		{Name: "subjs", Method: "go", Bin: "subjs",
			Install: goInstall("github.com/lc/subjs")},
		{Name: "LinkFinder", Method: "git", Bin: "linkfinder",
			Install: gitPy("https://github.com/GerbenJavado/LinkFinder.git", "LinkFinder", "linkfinder.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "SecretFinder", Method: "git", Bin: "secretfinder",
			Install: gitPy("https://github.com/m4ll0k/SecretFinder.git", "SecretFinder", "SecretFinder.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "JSluice", Method: "go", Bin: "jsluice",
			Install: goInstall("github.com/BishopFox/jsluice/cmd/jsluice")},
	}
}

// ── Fuzzing & content discovery ───────────────────────────────────────────

func fuzzTools() []Tool {
	return []Tool{
		{Name: "ffuf", Method: "go", Bin: "ffuf",
			Install: goInstall("github.com/ffuf/ffuf/v2")},
		{Name: "feroxbuster", Method: "cargo", Bin: "feroxbuster",
			Install: cargoTool("feroxbuster")},
		{Name: "gobuster", Method: "go", Bin: "gobuster",
			Install: goInstall("github.com/OJ/gobuster/v3")},
		{Name: "dirsearch", Method: "pipx", Bin: "dirsearch",
			Install: pipxTool("dirsearch")},
		{Name: "wfuzz", Method: "pip", Bin: "wfuzz",
			Install: pipUser("wfuzz")},
		{Name: "dirb", Method: "apt", Bin: "dirb", Install: apt("dirb")},
		{Name: "kiterunner", Method: "git", Bin: "kr",
			Install: func() error {
				dest := toolsDir + "/kiterunner"
				if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
					out, err := sh("git clone --depth=1 https://github.com/assetnote/kiterunner.git " + dest)
					if err != nil {
						return errFrom(out, err)
					}
				}
				out, err := sh("cd " + dest + " && make build")
				if err != nil {
					return errFrom(out, err)
				}
				findOut, _ := sh("find " + dest + " -maxdepth 3 -type f -name kr -perm -u+x 2>/dev/null | head -n1")
				found := strings.TrimSpace(findOut)
				if found == "" {
					return fmt.Errorf("build finished but 'kr' binary was not found under %s", dest)
				}
				linkDest := userBinDir + "/kr"
				os.Remove(linkDest)
				return os.Symlink(found, linkDest)
			},
			Info: "If the build fails, grab a prebuilt binary instead: https://github.com/assetnote/kiterunner/releases"},
	}
}

// ── Parameter discovery ───────────────────────────────────────────────────

func paramTools() []Tool {
	return []Tool{
		{Name: "arjun", Method: "pipx", Bin: "arjun", Install: pipxTool("arjun")},
		{Name: "x8", Method: "cargo", Bin: "x8", Install: cargoTool("x8")},
		{Name: "ParamSpider", Method: "git", Bin: "paramspider",
			Install: gitPy("https://github.com/devanshbatham/ParamSpider.git", "ParamSpider", "paramspider/main.py",
				"pip3 install --user .")},
		{Name: "unfurl", Method: "go", Bin: "unfurl", Install: goInstall("github.com/tomnomnom/unfurl")},
		{Name: "qsreplace", Method: "go", Bin: "qsreplace", Install: goInstall("github.com/tomnomnom/qsreplace")},
		{Name: "uro", Method: "pipx", Bin: "uro", Install: pipxTool("uro")},
	}
}

// ── Vulnerability scanning ────────────────────────────────────────────────

func vulnTools() []Tool {
	return []Tool{
		{Name: "nuclei", Method: "go", Bin: "nuclei",
			Install: goInstall("github.com/projectdiscovery/nuclei/v3/cmd/nuclei")},
		{Name: "nuclei-templates", Method: "git",
			CheckFunc: func() bool {
				_, err := os.Stat(homeDir() + "/nuclei-templates")
				return err == nil
			},
			Install: func() error {
				dest := homeDir() + "/nuclei-templates"
				if _, err := exec.LookPath("nuclei"); err == nil {
					if _, err := sh("nuclei -update-templates -silent"); err == nil {
						return nil
					}
				}
				out, err := sh("git clone --depth=1 https://github.com/projectdiscovery/nuclei-templates.git " + dest)
				return errFrom(out, err)
			}},
		{Name: "nikto", Method: "apt", Bin: "nikto", Install: apt("nikto")},
		{Name: "testssl.sh", Method: "git", Bin: "testssl.sh",
			Install: func() error {
				dest := toolsDir + "/testssl.sh"
				if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
					out, err := sh("git clone --depth=1 https://github.com/drwetter/testssl.sh.git " + dest)
					if err != nil {
						return errFrom(out, err)
					}
				}
				sh("chmod +x " + dest + "/testssl.sh")
				linkDest := userBinDir + "/testssl.sh"
				os.Remove(linkDest)
				return os.Symlink(dest+"/testssl.sh", linkDest)
			}},
	}
}
