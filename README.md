<a id="readme-top"></a>

<br />
<div align="center">
  <h3 align="center">MinecraftCrawler</h3>

  <p align="center">
    High-performance Minecraft server crawler and analyzer written in Go.
    <br />
    <br />
    <a href="#quick-start">Quick Start</a>
    ·
    <a href="#usage">Usage</a>
    ·
    <a href="#troubleshooting">Troubleshooting</a>
  </p>
</div>

```text
___  ____                            __ _     _____                    _
|  \/  (_)                          / _| |   /  __ \                  | |
| .  . |_ _ __   ___  ___ _ __ __ _| |_| |_  | /  \/_ __ __ ___      _| | ___ _ __
| |\/| | | '_ \ / _ \/ __| '__/ _` |  _| __| | |   | '__/ _` \ \ /\ / / |/ _ \ '__|
| |  | | | | | |  __/ (__| | | (_| | | | |_  | \__/\ | | (_| |\ V  V /| |  __/ |
\_|  |_/_|_| |_|\___|\___|_|  \__,_|_|  \__|  \____/_|  \__,_| \_/\_/ |_|\___|_|
```

[![build](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-build.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-build.yml)
[![test](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-test.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/lint.yml/badge.svg)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/lint.yml)
[![CodeQL](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/codeql.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/codeql.yml)

## Table of Contents

- [About](#about)
- [Quick Start](#quick-start)
- [Compatibility](#compatibility)
- [Architecture](#architecture)
- [Usage](#usage)
- [Data Collected](#data-collected)
- [Automated Releases](#automated-releases)
- [Troubleshooting](#troubleshooting)
- [Security and Responsible Use](#security-and-responsible-use)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## About

MinecraftCrawler discovers, analyzes, and stores Minecraft server metadata at scale. It combines `masscan` for fast host discovery with a concurrent protocol analyzer in Go.

Key features:

- Massive range scanning with worker-based analysis.
- Deep single-target analysis (`info`) with SRV-aware resolution.
- SLP + UDP Query extraction (version, players, software, plugins, map, MOTD).
- RCON status probing and secure-chat/whitelist signal extraction when available.
- SQLite persistence with WAL mode and schema migration support.
- Colorized CLI output and automatic run logging (`crawler.log`).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Quick Start

```sh
git clone https://github.com/miguerubsk/MinecraftCrawler.git
cd MinecraftCrawler
go build -o mccrawler main.go
./mccrawler info mc.hypixel.net
```

For mass range scanning (requires elevated privileges for `masscan` raw sockets):

```sh
sudo ./mccrawler scan --range 1.1.0.0/16 --rate 5000 --workers 2000
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Compatibility

- Go: `1.24+`
- OS: Linux, macOS, Windows
- Architectures: `amd64`, `arm64`, `arm`
- Scanner dependency: `masscan` (required for `scan`, not required for `info`)

Built with:

- [Go](https://go.dev/)
- [SQLite](https://www.sqlite.org/)
- [Masscan](https://github.com/robertdavidgraham/masscan)
- [Cobra](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Architecture

Pipeline overview for large scans:

```text
masscan -> ipChan -> worker pool -> AnalyzeServer -> resultChan -> SQLite batch writer
```

High-level command responsibilities:

- `scan`: discovers hosts in CIDR ranges, analyzes each host, batches persistence to SQLite.
- `info`: deep analysis of one target with automatic host/port resolution and SRV lookup.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

### `scan` command

Use `scan` for large CIDR ranges.

```sh
sudo ./mccrawler scan --range 1.1.0.0/16 --rate 5000 --workers 2000
```

Options:

| Flag | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--range` | `-r` | Required CIDR range (e.g. `1.1.1.0/24`) | `""` |
| `--rate` | `-p` | Packets per second for `masscan` | `1000` |
| `--port` |  | Target port for scan/analyze | `25565` |
| `--workers` | `-w` | Concurrent analysis workers | `1000` |
| `--verbose` | `-v` | Show found servers in stdout (optional line limit; no-value defaults to `500`) | `0` |
| `--exclude` |  | Exclusion file for IP ranges | `""` |
| `--output` | `-o` | SQLite output file | `results.db` |

### `info` command

Use `info` for deep single-target analysis.

```sh
./mccrawler info mc.hypixel.net
./mccrawler info 192.168.1.100:25565
./mccrawler info [2001:db8::1]:25565
```

What `info` does:

- Resolves host and port from target input.
- Tries `_minecraft._tcp` SRV for domain targets without explicit port.
- Runs SLP, UDP Query probe, and RCON status probe (`25575`).
- Prints enriched server details in a structured terminal report.

Example output:

```text
[*] Analyzing target: mc.example.net
[*] No port specified, looking up SRV record for mc.example.net...
[+] SRV found: play.example.net:25565

  SERVER ANALYSIS
  Server:                play.example.net:25565
  Version:               1.20.4
  Players:               12 / 200
  Software:              Paper
  Map:                   world
  RCON:                  Closed
  UDP Query:             CONNECTED
```

Note: `--output/-o` is intentionally `scan`-only and does not apply to `info`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Data Collected

Primary fields collected during analysis:

| Field | Source |
| --- | --- |
| `ip`, `port` | target resolution / scanner input |
| `version_name`, `protocol` | SLP |
| `players_online`, `players_max` | SLP |
| `motd` | SLP |
| `software`, `plugins`, `map_name` | UDP Query |
| `whitelist` | login/disconnect signal inference |
| `secure_chat` | SLP |
| `rcon_open` | RCON probe |
| `timestamp` | crawler runtime |

Persisted in SQLite by `scan`: core server metadata, software, map, mods/plugins JSON, whitelist, secure_chat, timestamp.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Automated Releases

Releases are automated through GitHub Actions (`release-please` + `goreleaser`).

Current release target matrix:

- OS: `linux`, `darwin`, `windows`
- Arch: `amd64`, `arm64`, `arm`
- Archives:
  - `tar.gz` for Linux/macOS
  - `zip` for Windows

Checksums are generated and signed in the release pipeline.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Troubleshooting

Common issues and fixes:

- `masscan` permission errors:
  - Run `scan` with elevated privileges (`sudo` on Linux/macOS).
  - `info` does not need raw-socket permissions.
- SRV not found:
  - This is normal for many servers; crawler falls back to default port `25565`.
- IPv6 target fails:
  - Use bracketed form when passing explicit port: `[2001:db8::1]:25565`.
- `Query UDP: Inactivo`:
  - Server may have Query disabled or filtered by firewall.
- `RCON: Cerrado`:
  - Expected for most servers unless RCON is intentionally exposed.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Security and Responsible Use

This tool is for educational and authorized assessment use only.

- Only scan infrastructure you own or explicitly have permission to test.
- Follow applicable laws and organizational policies.
- Avoid disruptive scan rates on sensitive networks.

The authors are not responsible for misuse.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Roadmap

Completed:

- [x] Massive scanning with Masscan
- [x] SLP protocol analysis
- [x] UDP Query extraction (software/plugins/map)
- [x] RCON status probe
- [x] Single-target deep analysis (`info`)
- [x] Optimized SQLite storage and migration-safe schema updates

Planned:

- [ ] Export to JSON/CSV
- [ ] Web dashboard
- [ ] GeoIP enrichment
- [ ] Webhook notifications
- [ ] Favicon extraction and local storage
- [ ] Search/filter command over stored results
- [ ] RCON password auditing module

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are welcome.

Local dev checks:

```sh
go test ./...
```

```sh
golangci-lint run
```

Typical flow:

1. Fork the project
2. Create a branch (`git checkout -b feature/my-change`)
3. Commit (`git commit -m "feat: ..."`)
4. Push and open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the GPL-3.0 License. See `LICENSE`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
