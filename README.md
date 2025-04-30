# 🕸️ ssrf-mapper

A powerful command-line tool for **automated SSRF fuzzing, metadata endpoint probing, and internal network discovery** — built in Go.

---

## 🚀 Why?

Server-Side Request Forgery (SSRF) continues to be one of the most dangerous vulnerabilities in modern cloud-native apps. While tools like `ffuf` and Burp Suite assist in testing, **there’s no dedicated CLI tool focused on exploring internal metadata, port scanning, redirect chaining, and evasion payloads — until now**.

---

## ✨ Features

- 🔁 **Fuzz with wordlist-based parameters**
- 🧠 **Auto-detect cloud provider** (AWS, GCP, Azure)
- 🎯 **Test multiple SSRF encodings**: raw, URL-encoded, decimal, hex, octal, IPv6-mapped
- 🔐 **Send custom headers & cookies**
- 🛰️ **Probe metadata services** like `169.254.169.254`
- 🛣️ **Handle redirect chains** (30x responses)
- 🌐 **Internal port scanning** (on common ranges)
- 📤 **JSON/CSV logging for reports**
- ⚡ **Highly concurrent & fast**
- **JS Redirect Detection**: Automatically detects JavaScript-based client-side redirects (useful for detecting sites that redirect via JavaScript).
- **CLI-based**: Simple command-line interface for scanning URLs.
- **Configurable**: Flags to enable or disable features like JS redirect detection.

---

## 📦 Installation

```bash
git clone https://github.com/jayasurya963/ssrf-mapper.git
cd ssrf-mapper
go build -o ssrf-mapper ./cmd/ssrf-mapper
