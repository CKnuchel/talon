# Talon
![Go Version](https://img.shields.io/github/go-mod/go-version/CKnuchel/talon)
![Release](https://img.shields.io/github/v/tag/CKnuchel/talon?label=release)

![Tests](https://img.shields.io/github/actions/workflow/status/CKnuchel/talon/test.yml?branch=master&label=tests)
![License](https://img.shields.io/github/license/CKnuchel/talon)



A lightweight, growing CLI toolkit for network reconnaissance, inspired by nmap. Built as a hands-on Go learning project.

## Status

Early development. Currently supports concurrent TCP port scanning with structured output.

## Features

- Concurrent TCP port scanning with a configurable worker pool
- IPv4 and IPv6 support
- Table output in the terminal
- Export results as JSON or CSV

## Installation

Download a prebuilt binary from the [Releases](https://github.com/CKnuchel/talon/releases) page, or build from source:

```bash
git clone https://github.com/CKnuchel/talon.git
cd talon
go build -o talon ./cmd/talon
```

## Usage

```bash
talon [flags] <host>
```

### Flags

| Flag | Description                          | Default |
| ---- | ------------------------------------- | ------- |
| `-s` | Start port                           | `1`     |
| `-e` | End port                             | `1024`  |
| `-t` | Connection timeout                   | `200ms` |
| `-w` | Number of parallel scan workers      | `100`   |
| `-c` | Include closed ports in the result   | `false` |
| `-o` | Export results to a file (`.json` or `.csv`) | — |

### Examples

Scan the default port range:
```bash
talon scanme.nmap.org
```

Scan a custom range with more workers:
```bash
talon -s 1 -e 65535 -w 300 scanme.nmap.org
```

Export results as JSON:
```bash
talon -s 1 -e 1024 -o results.json scanme.nmap.org
```

## License

[MIT](LICENSE)