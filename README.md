<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h3 align="center">MinecraftCrawler</h3>

  <p align="center">
    An ultra-high performance Minecraft server crawler written in Go.
    <br />
    <br />
    <a href="#usage">View Demo</a>
    ·
    <a href="#issues">Report Bug</a>
    ·
    <a href="#issues">Request Feature</a>
  </p>
</div>

[![build](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-build.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-build.yml)
[![test](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-test.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/go-test.yml)
[![golangci-lint](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/lint.yml/badge.svg)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/lint.yml)
[![CodeQL](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/codeql.yml/badge.svg?branch=master)](https://github.com/miguerubsk/MinecraftCrawler/actions/workflows/codeql.yml)

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

MinecraftCrawler is a tool designed to discover, analyze, and store Minecraft server information at a massive scale. It combines the port discovery speed of `masscan` with a highly concurrent protocol analyzer (SLP) written in Go.

Key Features:

- **Extreme Speed**: Pipeline architecture capable of processing thousands of servers per second.
- **Efficiency**: Optimized use of goroutines and SQLite database with WAL mode for batch writing.
- **Deep Analysis**: Extracts version, players, MOTD, mod list (Forge), plugins, map name, whitelist status, UDP Query status, and RCON probe status.
- **Premium CLI**: Highly professional terminal interface with ASCII branding, color-coded output, and real-time feedback.
- **Automatic Logging**: Every session is automatically recorded in `crawler.log` for later analysis.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

- [![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?logo=go&logoColor=white)](https://go.dev/)
- [![SQLite](https://img.shields.io/badge/SQLite-%2307405E.svg?logo=sqlite&logoColor=white)](https://www.sqlite.org/)
- [Masscan](https://github.com/robertdavidgraham/masscan)
- [Cobra](https://github.com/spf13/cobra)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

You need to have Go and Masscan installed on your system (Linux/Windows/Mac).

- **Masscan** (Debian/Ubuntu):

  ```sh
  sudo apt-get install masscan
  ```

- **Go** (1.24+):
  Download it from [go.dev/dl](https://go.dev/dl/).

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/miguerubsk/MinecraftCrawler.git
   ```
2. Install Go dependencies
   ```sh
   go mod tidy
   ```
3. Build the binary
   ```sh
   go build -o mccrawler main.go
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

You can use either `scan` for large ranges or `info` for deep single-target analysis.

### `scan` command

You need administrator privileges to run `masscan` (network raw sockets).

**Basic Example:** Scan a full IP range

```sh
sudo ./mccrawler scan --range 1.1.0.0/16 --rate 5000 --workers 2000
```

**Available Options:**

| Flag        | Shorthand | Description                                               | Default      |
| ----------- | --------- | --------------------------------------------------------- | ------------ |
| `--range`   | `-r`      | **Required**. CIDR range to scan (e.g., 1.1.1.0/24)       | `""`         |
| `--rate`    | `-p`      | Packets per second (Masscan)                              | `1000`       |
| `--port`    |           | Port to scan (25565 or 25575)                             | `25565`      |
| `--workers` | `-w`      | Number of concurrent worker threads                       | `1000`       |
| `--verbose` | `-v`      | Display found servers (limit lines, optional default 500) | `0`          |
| `--exclude` |           | IP exclusion file                                         | `""`         |
| `--output`  | `-o`      | Output database file                                      | `results.db` |

_Check `mccrawler help` for more information. All execution logs are automatically saved to `crawler.log`. At the end of each scan, a statistical summary dashboard is displayed._

### `info` command

Run an in-depth analysis for one target (domain/IP), including:

- SRV auto-resolution (`_minecraft._tcp`)
- Server List Ping (SLP)
- UDP Query probe and parsing (software/plugins/map)
- RCON status probe (`25575`)
- Whitelist and secure chat signals when available

Examples:

```sh
./mccrawler info mc.hypixel.net
./mccrawler info 192.168.1.100:25565
./mccrawler info [2001:db8::1]:25565
```

Note: `--output/-o` is a `scan`-only flag. It does not apply to `info`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [x] Massive scanning with Masscan
- [x] SLP (Server List Ping) protocol analysis
- [x] Whitelist and Mods detection
- [x] Optimized SQLite storage
- [x] RCON status probe support
- [ ] Export to JSON/CSV format
- [ ] Web dashboard for result visualization
- [x] **Query Protocol (UDP)** support for software/plugins/map extraction
- [ ] **GeoIP Integration** (Country/Flag detection)
- [ ] **Webhook Notifications** (Discord/Slack support)
- [ ] **Favicon Extraction** and local PNG storage
- [ ] **CLI Search** command for quick result filtering
- [ ] **RCON Password Auditing** module
- [x] **Single Target Scan** (`info`) with auto-detect target resolution

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Legal Disclaimer

This tool is for educational purposes only. Unauthorized scanning of networks without permission can be illegal and unethical. Always ensure you have proper authorization before using this tool. The author is not responsible for its misuse or the consequences of scanning networks without authorization.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
