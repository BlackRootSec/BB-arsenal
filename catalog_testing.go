package main

import "os"

// ── XSS testing ────────────────────────────────────────────────────────────

func xssTools() []Tool {
	return []Tool{
		{Name: "dalfox", Method: "go", Bin: "dalfox", Install: goInstall("github.com/hahwul/dalfox/v2")},
		{Name: "kxss", Method: "go", Bin: "kxss", Install: goInstall("github.com/Emoe/kxss")},
		{Name: "Gxss", Method: "go", Bin: "Gxss", Install: goInstall("github.com/KathanP19/Gxss")},
		{Name: "XSStrike", Method: "git", Bin: "xsstrike",
			Install: gitPy("https://github.com/s0md3v/XSStrike.git", "XSStrike", "xsstrike.py",
				"pip3 install --user -r requirements.txt")},
	}
}

// ── SQL injection ──────────────────────────────────────────────────────────

func sqliTools() []Tool {
	return []Tool{
		{Name: "sqlmap", Method: "apt", Bin: "sqlmap", Install: apt("sqlmap")},
		{Name: "ghauri", Method: "pipx", Bin: "ghauri", Install: pipxTool("ghauri")},
	}
}

// ── API, proxy & protocol tooling ──────────────────────────────────────────

func apiTools() []Tool {
	return []Tool{
		{Name: "httpie", Method: "apt", Bin: "http", Install: apt("httpie")},
		{Name: "mitmproxy", Method: "pipx", Bin: "mitmproxy", Install: pipxTool("mitmproxy")},
		{Name: "interactsh-client", Method: "go", Bin: "interactsh-client",
			Install: goInstall("github.com/projectdiscovery/interactsh/cmd/interactsh-client")},
		{Name: "InQL", Method: "pipx", Bin: "inql", Install: pipxTool("inql"),
			Info: "Standalone CLI installed. A Burp Suite version is also available via the BApp Store."},
	}
}

// ── JWT tooling ─────────────────────────────────────────────────────────────

func jwtTools() []Tool {
	return []Tool{
		{Name: "jwt_tool", Method: "git", Bin: "jwt_tool",
			Install: gitPy("https://github.com/ticarpi/jwt_tool.git", "jwt_tool", "jwt_tool.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "PyJWT", Method: "pip", CheckFunc: pyModuleInstalled("jwt"), Install: pipUser("pyjwt"),
			Info: "Python library — import with `import jwt` in scripts."},
	}
}

// ── Secrets & static analysis ──────────────────────────────────────────────

func secretsTools() []Tool {
	return []Tool{
		{Name: "trufflehog", Method: "script", Bin: "trufflehog",
			Install: rawCmd("curl -sSfL https://raw.githubusercontent.com/trufflesecurity/trufflehog/main/scripts/install.sh | " + sudoCmd + "sh -s -- -b /usr/local/bin")},
		{Name: "gitleaks", Method: "go", Bin: "gitleaks", Install: goInstall("github.com/gitleaks/gitleaks/v8")},
		{Name: "git-secrets", Method: "git", Bin: "git-secrets",
			Install: func() error {
				dest := toolsDir + "/git-secrets"
				if _, statErr := os.Stat(dest); os.IsNotExist(statErr) {
					out, err := sh("git clone --depth=1 https://github.com/awslabs/git-secrets.git " + dest)
					if err != nil {
						return errFrom(out, err)
					}
				}
				out, err := sh("cd " + dest + " && " + sudoCmd + "make install")
				return errFrom(out, err)
			}},
		{Name: "detect-secrets", Method: "pipx", Bin: "detect-secrets", Install: pipxTool("detect-secrets")},
		{Name: "semgrep", Method: "pipx", Bin: "semgrep", Install: pipxTool("semgrep")},
		{Name: "Retire.js", Method: "npm", Bin: "retire", Install: npmGlobal("retire")},
	}
}

// ── Tech fingerprinting ─────────────────────────────────────────────────────

func techTools() []Tool {
	return []Tool{
		{Name: "whatweb", Method: "apt", Bin: "whatweb", Install: apt("whatweb")},
		{Name: "webanalyze", Method: "go", Bin: "webanalyze",
			Install: func() error {
				out, err := sh("go install -v github.com/rverton/webanalyze/cmd/webanalyze@latest")
				if err != nil {
					return errFrom(out, err)
				}
				sh("webanalyze -update")
				return nil
			},
			Info: "Uses the Wappalyzer fingerprint database — covers the same detections as the Wappalyzer browser extension."},
	}
}

// ── Cloud security ──────────────────────────────────────────────────────────

func cloudTools() []Tool {
	return []Tool{
		{Name: "cloud_enum", Method: "git", Bin: "cloud_enum",
			Install: gitPy("https://github.com/initstring/cloud_enum.git", "cloud_enum", "cloud_enum.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "CloudBrute", Method: "go", CheckFunc: anyBin("CloudBrute", "cloudbrute"),
			Install: goInstall("github.com/0xsha/CloudBrute")},
		{Name: "Prowler", Method: "pipx", Bin: "prowler", Install: pipxTool("prowler")},
		{Name: "ScoutSuite", Method: "pipx", Bin: "scout", Install: pipxTool("scoutsuite")},
		{Name: "Pacu", Method: "git", Bin: "pacu",
			Install: gitPy("https://github.com/RhinoSecurityLabs/pacu.git", "Pacu", "pacu.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "S3Scanner", Method: "go", CheckFunc: anyBin("S3Scanner", "s3scanner"),
			Install: goInstall("github.com/sa7mon/S3Scanner")},
	}
}

// ── Mobile & reverse engineering ────────────────────────────────────────────

func mobileTools() []Tool {
	return []Tool{
		{Name: "apktool", Method: "apt", Bin: "apktool", Install: apt("apktool")},
		{Name: "adb", Method: "apt", Bin: "adb", Install: apt("android-tools-adb")},
		{Name: "jadx", Method: "apt", Bin: "jadx", Install: apt("jadx")},
		{Name: "frida-tools", Method: "pipx", Bin: "frida", Install: pipxTool("frida-tools")},
		{Name: "objection", Method: "pipx", Bin: "objection", Install: pipxTool("objection")},
		{Name: "radare2", Method: "apt", Bin: "r2", Install: apt("radare2")},
		{Name: "apk-mitm", Method: "npm", Bin: "apk-mitm", Install: npmGlobal("apk-mitm")},
		{Name: "binutils (strings/objdump/readelf/nm)", Method: "apt", Bin: "objdump", Install: apt("binutils")},
	}
}

// ── Network utilities ────────────────────────────────────────────────────────

func netUtilTools() []Tool {
	return []Tool{
		{Name: "socat", Method: "apt", Bin: "socat", Install: apt("socat")},
		{Name: "netcat-openbsd", Method: "apt", Bin: "nc", Install: apt("netcat-openbsd")},
		{Name: "tcpdump", Method: "apt", Bin: "tcpdump", Install: apt("tcpdump")},
		{Name: "tshark", Method: "apt", Bin: "tshark", Install: apt("tshark")},
		{Name: "dnsutils", Method: "apt", Bin: "dig", Install: apt("dnsutils")},
		{Name: "whois", Method: "apt", Bin: "whois", Install: apt("whois")},
		{Name: "net-tools", Method: "apt", Bin: "netstat", Install: apt("net-tools")},
		{Name: "proxychains4", Method: "apt", Bin: "proxychains4", Install: apt("proxychains4")},
	}
}

// ── OSINT ────────────────────────────────────────────────────────────────────

func osintTools() []Tool {
	return []Tool{
		{Name: "theHarvester", Method: "apt", Bin: "theHarvester", Install: apt("theharvester")},
		{Name: "SpiderFoot", Method: "git", Bin: "spiderfoot",
			Install: gitPy("https://github.com/smicallef/spiderfoot.git", "SpiderFoot", "sf.py",
				"pip3 install --user -r requirements.txt")},
		{Name: "Shodan CLI", Method: "pipx", Bin: "shodan", Install: pipxTool("shodan"),
			Info: "Run `shodan init <API_KEY>` once you have a Shodan account."},
		{Name: "Censys CLI", Method: "pipx", Bin: "censys", Install: pipxTool("censys"),
			Info: "Run `censys config` once you have a Censys account."},
	}
}

// ── Optional GUI apps (skippable with --no-gui) ────────────────────────────

func guiTools() []Tool {
	return []Tool{
		{Name: "Burp Suite Community", Method: "apt", Bin: "burpsuite", Install: apt("burpsuite"),
			Info: "Not packaged on plain Ubuntu/Debian — install manually from https://portswigger.net/burp if apt can't find it."},
		{Name: "Wireshark", Method: "apt", Bin: "wireshark", Install: apt("wireshark")},
		{Name: "Postman", Method: "snap", Bin: "postman", Install: rawCmd(sudoCmd + "snap install postman")},
		{Name: "Insomnia", Method: "snap", Bin: "insomnia", Install: rawCmd(sudoCmd + "snap install insomnia")},
	}
}
