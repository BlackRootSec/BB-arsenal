package main

// ── Category assembly ───────────────────────────────────────────────────────

func buildCategories() []Category {
	return []Category{
		{Key: "prereq", Title: "Prerequisites", Icon: "🧰", Tools: prereqTools()},
		{Key: "subs", Title: "Subdomain Enumeration", Icon: "🌐", Tools: subdomainTools()},
		{Key: "dns", Title: "DNS Resolution & Permutation", Icon: "🧭", Tools: dnsTools()},
		{Key: "http", Title: "HTTP Probing & Port Scanning", Icon: "🚪", Tools: httpPortTools()},
		{Key: "crawl", Title: "Crawling & URL Discovery", Icon: "🕷️", Tools: crawlTools()},
		{Key: "fuzz", Title: "Fuzzing & Content Discovery", Icon: "🎯", Tools: fuzzTools()},
		{Key: "params", Title: "Parameter Discovery", Icon: "🔧", Tools: paramTools()},
		{Key: "vuln", Title: "Vulnerability Scanning", Icon: "🛡️", Tools: vulnTools()},
		{Key: "xss", Title: "XSS Testing", Icon: "💉", Tools: xssTools()},
		{Key: "sqli", Title: "SQL Injection", Icon: "🗄️", Tools: sqliTools()},
		{Key: "api", Title: "API, Proxy & Protocol Tools", Icon: "🔌", Tools: apiTools()},
		{Key: "jwt", Title: "JWT Tooling", Icon: "🔑", Tools: jwtTools()},
		{Key: "secrets", Title: "Secrets & Static Analysis", Icon: "🕵️", Tools: secretsTools()},
		{Key: "tech", Title: "Tech Fingerprinting", Icon: "🔍", Tools: techTools()},
		{Key: "cloud", Title: "Cloud Security", Icon: "☁️", Tools: cloudTools()},
		{Key: "mobile", Title: "Mobile & Reverse Engineering", Icon: "📱", Tools: mobileTools()},
		{Key: "netutils", Title: "Network Utilities", Icon: "🛰️", Tools: netUtilTools()},
		{Key: "osint", Title: "OSINT", Icon: "🛰️", Tools: osintTools()},
		{Key: "gui", Title: "Optional GUI Apps", Icon: "🖥️", Tools: guiTools()},
	}
}
