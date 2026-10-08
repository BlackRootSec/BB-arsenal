package main

// ManualNote documents something from the original wishlist that this
// installer intentionally does not try to "install", and why.
type ManualNote struct {
	Name   string
	Reason string
}

func manualReconTools() []ManualNote {
	return []ManualNote{
		{"Ghidra", "GUI reverse-engineering suite that needs a JDK and a large manual download. Get it from https://ghidra-sre.org — no reliable silent-install path."},
		{"Cutter", "GUI front-end for radare2, distributed as an AppImage. Download from https://cutter.re — radare2/r2 (already installed above) covers the CLI side."},
		{"MobSF", "Best run via Docker: `docker pull opensecurity/mobile-security-framework-mobsf` then `docker run -it -p 8000:8000 opensecurity/mobile-security-framework-mobsf`. Needs Docker installed and some setup, so it's left manual."},
		{"x64dbg / WinDbg", "Windows-only debuggers — not applicable on a Linux recon box. On Linux, radare2/r2 or Ghidra above fill the same role."},
	}
}

func manualBurpExtensions() []ManualNote {
	return []ManualNote{
		{"DOM Invader", "Built into Burp Suite Pro — enable it from Burp's browser settings, nothing to install separately."},
		{"Autorize", "Burp extension for authorization testing — install from Burp's BApp Store (Extender tab)."},
		{"AuthMatrix", "Burp extension for role-based access testing — install from the BApp Store."},
		{"JWT Editor", "Burp extension for editing/forging JWTs — install from the BApp Store. jwt_tool and PyJWT (installed above) cover the CLI side."},
		{"Param Miner", "Burp extension for hidden parameter discovery — install from the BApp Store."},
		{"HTTP Request Smuggler", "Burp extension by James Kettle — install from the BApp Store."},
		{"GraphQL Voyager", "Browser-based GraphQL schema visualizer — run via `npx graphql-voyager` or its hosted demo; not a system package."},
	}
}

func manualWebServices() []ManualNote {
	return []ManualNote{
		{"crt.sh", "A certificate-transparency search website, not a CLI tool. subfinder/amass already query crt.sh as one of their passive sources."},
		{"jwt.io", "A website for decoding/debugging JWTs in the browser. jwt_tool and PyJWT (installed above) give you the same thing from the CLI."},
		{"Shodan / Censys", "Web services — CLI clients ARE installed above (`shodan`, `censys`), but you still need your own API key/account."},
		{"FOFA / ZoomEye / SecurityTrails", "Web services with paid APIs and no single standard CLI client. Sign up and use their web UI or REST API directly."},
	}
}

func manualConcepts() []ManualNote {
	return []ManualNote{
		{"REST / GraphQL / WebSockets / gRPC / SOAP / OpenAPI / Swagger", "These are protocols/specs, not software. Test them with the HTTP/API tools already installed above (curl, httpie, mitmproxy, InQL for GraphQL, Postman/Insomnia for REST)."},
		{"OAuth/OIDC, JWT, API keys, BOLA/IDOR, BFLA, mass assignment, rate limiting, race conditions", "These are vulnerability classes/concepts to test for, not installable tools. jwt_tool/PyJWT target JWT issues; Autorize/AuthMatrix (Burp BApp Store, noted above) target BOLA/IDOR/BFLA; ffuf/dalfox/etc. are your general-purpose testing tools for the rest."},
	}
}
